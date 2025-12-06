package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/internal/service"

	"github.com/gin-gonic/gin"
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

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.FailWithMessage("未上传文件", c)
		return
	}
	defer file.Close()

	doc, err := documentService.Upload(uint(libraryID), file, header)
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
	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// 获取 flusher
	flusher, ok := c.Writer.(interface{ Flush() })
	if !ok {
		c.SSEvent("error", "SSE not supported")
		return
	}

	// 发送 SSE 事件的辅助函数
	sendEvent := func(eventType string, data interface{}) {
		jsonData, _ := json.Marshal(data)
		fmt.Fprintf(c.Writer, "data: %s\n\n", jsonData)
		flusher.Flush()
	}

	// 解析参数
	libraryIDStr := c.PostForm("library_id")
	libraryID, err := strconv.ParseUint(libraryIDStr, 10, 32)
	if err != nil {
		sendEvent("error", map[string]string{"message": "无效的库ID"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		sendEvent("error", map[string]string{"message": "未上传文件"})
		return
	}
	defer file.Close()

	// 创建状态通道
	statusChan := make(chan response.ProcessStatus, 10)

	// 上传文档（带状态回调）
	doc, err := documentService.UploadWithCallback(uint(libraryID), file, header, statusChan)
	if err != nil {
		errMsg := "上传失败"
		if errors.Is(err, service.ErrNotFound) {
			errMsg = "库不存在"
		} else if errors.Is(err, service.ErrAlreadyExists) {
			errMsg = "相同内容的文档已存在"
		} else if errors.Is(err, service.ErrInvalidParams) {
			errMsg = "不支持的文件类型"
		}
		sendEvent("error", map[string]string{"message": errMsg})
		return
	}

	// 发送上传成功事件
	sendEvent("uploaded", map[string]interface{}{
		"document_id": doc.ID,
		"title":       doc.Title,
		"status":      doc.Status,
	})

	// 监听处理状态
	for status := range statusChan {
		sendEvent(status.Stage, map[string]interface{}{
			"document_id": doc.ID,
			"stage":       status.Stage,
			"progress":    status.Progress,
			"message":     status.Message,
			"status":      status.Status,
		})

		// 处理完成或失败，退出
		if status.Stage == "completed" || status.Stage == "failed" {
			break
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

// GetLatestCode 获取库的最新文档内容
// GET /documents/code/:libid
func (d *DocumentApi) GetLatestCode(c *gin.Context) {
	libraryID, err := strconv.ParseUint(c.Param("libid"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的库ID", c)
		return
	}

	title, content, err := documentService.GetLatestContent(uint(libraryID))
	if err != nil {
		response.OkWithData(gin.H{
			"title":   "",
			"content": "No documents available. Upload a document to get started.",
		}, c)
		return
	}

	response.OkWithData(gin.H{
		"title":   title,
		"content": content,
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
