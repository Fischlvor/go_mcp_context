package api

import (
	"errors"
	"strconv"

	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/internal/service"
	"go-mcp-context/pkg/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DocumentApi struct{}

// List 获取文档列表
func (d *DocumentApi) List(c *gin.Context) {
	var req request.DocumentList
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	result, err := documentService.List(&req)
	if err != nil {
		response.FailWithMessage("查询失败: "+err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// Upload 上传文档
func (d *DocumentApi) Upload(c *gin.Context) {
	libraryIDStr := c.PostForm("library_id")
	libraryID, err := strconv.ParseUint(libraryIDStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的库ID", c)
		return
	}

	version := c.PostForm("version")
	if version == "" {
		version = "latest" // 默认版本
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.FailWithMessage("未上传文件", c)
		return
	}
	defer file.Close()

	doc, err := documentService.Upload(uint(libraryID), version, file, header)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.FailWithMessage("库不存在", c)
			return
		}
		if errors.Is(err, service.ErrAlreadyExists) {
			response.FailWithMessage("相同内容的文档已存在", c)
			return
		}
		if errors.Is(err, service.ErrInvalidParams) {
			response.FailWithMessage("不支持的文件类型", c)
			return
		}
		response.FailWithMessage("上传失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(doc, "上传成功，正在处理中", c)
}

// UploadWithSSE 上传文档（SSE 实时推送处理状态）
func (d *DocumentApi) UploadWithSSE(c *gin.Context) {
	// 创建文档上传专用 SSE 写入器
	sse, ok := response.NewDocumentSSEWriter(c)
	if !ok {
		c.SSEvent("error", "SSE not supported")
		return
	}

	// 解析参数
	libraryIDStr := c.PostForm("library_id")
	libraryID, err := strconv.ParseUint(libraryIDStr, 10, 32)
	if err != nil {
		sse.SendError("无效的库ID")
		return
	}

	version := c.PostForm("version")
	if version == "" {
		version = "latest" // 默认版本
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		sse.SendError("未上传文件")
		return
	}
	defer file.Close()

	// 创建状态通道
	statusChan := make(chan response.ProcessStatus, 10)

	// 上传文档（带状态回调）
	doc, err := documentService.UploadWithCallback(uint(libraryID), version, file, header, statusChan)
	if err != nil {
		errMsg := "上传失败"
		if errors.Is(err, service.ErrNotFound) {
			errMsg = "库不存在"
		} else if errors.Is(err, service.ErrAlreadyExists) {
			errMsg = "相同内容的文档已存在"
		} else if errors.Is(err, service.ErrInvalidParams) {
			errMsg = "不支持的文件类型"
		}
		global.Log.Error("Document upload failed", zap.String("error", errMsg), zap.Error(err))
		sse.SendError(errMsg)
		return
	}

	// 发送上传成功事件
	sse.SendProgress("uploaded", 5, "文件上传成功", doc.ID)

	// 监听处理状态
	for status := range statusChan {
		if status.Stage == "completed" {
			sse.SendComplete(doc.ID, doc.Title)
			break
		} else if status.Stage == "failed" {
			sse.SendFailed(status.Message, doc.ID)
			break
		} else {
			sse.SendProgress(status.Stage, status.Progress, status.Message, doc.ID)
		}
	}
}

// Get 获取文档详情
func (d *DocumentApi) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的文档ID", c)
		return
	}

	doc, err := documentService.GetByID(uint(id))
	if err != nil {
		response.FailWithMessage("文档不存在", c)
		return
	}

	response.OkWithData(doc, c)
}

// GetChunks 获取库的文档块
// GET /documents/:mode/:libid/*version
// mode: code 或 info, version 可选
func (d *DocumentApi) GetChunks(c *gin.Context) {
	mode := c.Param("mode") // code 或 info
	if mode != "code" && mode != "info" {
		response.FailWithMessage("无效的模式，必须是 code 或 info", c)
		return
	}

	libraryID, err := strconv.ParseUint(c.Param("libid"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的库ID", c)
		return
	}

	// 处理版本参数（*version 会带有前导斜杠）
	version := c.Param("version")
	if len(version) > 0 && version[0] == '/' {
		version = version[1:]
	}

	// 如果没有指定版本，使用库的默认版本
	if version == "" {
		library, err := libraryService.GetByID(uint(libraryID))
		if err != nil {
			response.FailWithMessage("库不存在", c)
			return
		}
		version = library.DefaultVersion
	}

	chunks, err := documentService.GetChunks(uint(libraryID), version, mode)
	if err != nil {
		response.OkWithData(gin.H{
			"chunks": []interface{}{},
		}, c)
		return
	}

	response.OkWithData(gin.H{
		"chunks": chunks,
	}, c)
}

// Delete 删除文档
func (d *DocumentApi) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的文档ID", c)
		return
	}

	if err := documentService.Delete(uint(id)); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.FailWithMessage("文档不存在", c)
			return
		}
		response.FailWithMessage("删除失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)
}
