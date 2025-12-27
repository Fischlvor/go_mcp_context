package api

import (
	"errors"
	"fmt"
	"strconv"

	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/internal/service"
	"go-mcp-context/pkg/bufferedwriter/actlog"
	"go-mcp-context/pkg/global"
	"go-mcp-context/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DocumentApi struct{}

// List 获取文档列表
// @Summary 获取文档列表
// @Description 分页获取指定库的文档列表
// @Tags Documents
// @Accept json
// @Produce json
// @Param library_id query int true "库 ID"
// @Param version query string false "版本号（可选）"
// @Param page query int false "页码，默认 1" default(1)
// @Param limit query int false "每页数量，默认 10，最大 50" default(10)
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 400 {object} response.Response
// @Router /api/v1/documents/list [get]
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
// @Summary 上传文档
// @Description 上传文档文件到指定库版本（支持 Markdown、PDF、DOCX、Swagger）（需要认证）
// @Tags Documents
// @Accept multipart/form-data
// @Produce json
// @Security JWTAuth
// @Param library_id formData int true "库 ID"
// @Param version formData string false "版本号，默认为 latest"
// @Param file formData file true "文档文件"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/documents/upload [post]
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

	// 获取用户和任务 ID（用于日志关联）
	userUUID := utils.GetUUID(c).String()
	taskID := utils.GenerateTaskID()

	doc, err := documentService.Upload(uint(libraryID), version, file, header, userUUID, taskID)
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

	// 活动日志已在 service 层记录
	response.OkWithDetailed(doc, "上传成功，正在处理中", c)
}

// UploadWithSSE 上传文档（SSE 实时推送处理状态）
// @Summary 上传文档（SSE 实时推送）
// @Description 上传文档文件到指定库版本，通过 SSE 实时推送处理进度（需要认证）
// @Tags Documents
// @Accept multipart/form-data
// @Produce text/event-stream
// @Security JWTAuth
// @Param library_id formData int true "库 ID"
// @Param version formData string false "版本号，默认为 latest"
// @Param file formData file true "文档文件"
// @Success 200 {object} response.ProcessStatus "处理状态流"
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/documents/upload-sse [post]
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
// @Summary 获取文档详情
// @Description 获取指定文档的完整信息
// @Tags Documents
// @Accept json
// @Produce json
// @Param id path int true "文档 ID"
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 404 {object} response.Response
// @Router /api/v1/documents/detail/:id [get]
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

// GetChunks 获取库的文档块（统一入口，支持搜索和列表）
// @Summary 获取文档块
// @Description 获取指定库的文档块，支持向量搜索和全量列表。mode 为 code 或 info，version 可选（默认使用库的默认版本），topic 可选（传入则进行向量搜索）
// @Tags Documents
// @Accept json
// @Produce json
// @Param mode path string true "模式：code(代码示例) 或 info(文档说明)"
// @Param libid path int true "库 ID"
// @Param version query string false "版本号，不传则使用库的默认版本"
// @Param topic query string false "搜索主题，不传则返回全部文档块"
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/documents/chunks/{mode}/{libid} [get]
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

	// 从 query 参数获取版本
	version := c.Query("version")

	// 如果没有指定版本，使用库的默认版本
	if version == "" {
		library, err := libraryService.GetByID(uint(libraryID))
		if err != nil {
			response.FailWithMessage("库不存在", c)
			return
		}
		version = library.DefaultVersion
	}

	// 获取 topic 参数（可选）
	topic := c.Query("topic")

	limit := 10

	// 如果有 topic，进行向量搜索
	if topic != "" {
		global.Log.Info("GetChunks: 执行向量搜索",
			zap.Uint64("library_id", libraryID),
			zap.String("mode", mode),
			zap.String("topic", topic),
		)

		// 调用搜索服务
		searchResult, err := searchService.SearchDocuments(&request.Search{
			LibraryID: uint(libraryID),
			Query:     topic,
			Mode:      mode,
			Version:   version, // 传入版本参数
			Page:      1,
			Limit:     limit, // 前端详情页返回 10 条
		})
		if err != nil {
			global.Log.Error("GetChunks: 搜索失败", zap.Error(err))
			response.FailWithMessage("搜索失败: "+err.Error(), c)
			return
		}

		// 将搜索结果转换为 chunks 格式
		chunks := make([]gin.H, len(searchResult.Results))
		for i, r := range searchResult.Results {
			chunks[i] = gin.H{
				"id":          r.ChunkID,
				"library_id":  r.LibraryID,
				"upload_id":   r.UploadID,
				"version":     r.Version,
				"title":       r.Title,
				"description": r.Description, // code mode: LLM 生成, info mode: 空
				"source":      r.Source,
				"language":    r.Language, // code mode: 代码语言, info mode: 空
				"code":        r.Code,     // code mode: 代码内容, info mode: 空
				"chunk_text":  r.Content,  // 原文内容
				"tokens":      r.Tokens,
				"chunk_type":  mode,
				"relevance":   r.Relevance,
			}
		}

		response.OkWithData(gin.H{
			"chunks": chunks,
			"total":  searchResult.Total,
			"topic":  topic,
		}, c)
		return
	}

	// 无 topic，返回文档块（受 limit 限制）
	dbChunks, err := documentService.GetChunks(uint(libraryID), version, mode, limit)
	if err != nil {
		response.OkWithData(gin.H{
			"chunks": []interface{}{},
		}, c)
		return
	}

	// 转换为最小所需字段格式
	chunks := make([]gin.H, len(dbChunks))
	for i, chunk := range dbChunks {
		chunks[i] = gin.H{
			"id":          chunk.ID,
			"title":       chunk.Title,
			"description": chunk.Description, // code mode: LLM 生成, info mode: 空
			"source":      chunk.Source,
			"language":    chunk.Language,  // code mode: 代码语言, info mode: 空
			"code":        chunk.Code,      // code mode: 代码内容, info mode: 空
			"chunk_text":  chunk.ChunkText, // 原文内容
			"tokens":      chunk.Tokens,
		}
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

	// 先获取文档信息（用于日志）
	doc, _ := documentService.GetByID(uint(id))

	if err := documentService.Delete(uint(id)); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.FailWithMessage("文档不存在", c)
			return
		}
		response.FailWithMessage("删除失败", c)
		return
	}

	// 记录活动日志
	if doc != nil {
		actlog.Success(doc.LibraryID, actlog.EventDocDelete, fmt.Sprintf("删除文档: %s", doc.Title),
			actlog.WithTarget("document", strconv.FormatUint(uint64(id), 10)),
			actlog.WithVersion(doc.Version))
	}

	response.OkWithMessage("删除成功", c)
}
