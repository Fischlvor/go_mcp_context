package api

import (
	"errors"
	"strconv"

	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/internal/service"

	"github.com/gin-gonic/gin"
)

type LibraryApi struct{}

// List 获取库列表（带统计信息）
func (l *LibraryApi) List(c *gin.Context) {
	var req request.LibraryList
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	result, err := libraryService.ListWithStats(&req)
	if err != nil {
		response.FailWithMessage("查询失败: "+err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// Create 创建库
func (l *LibraryApi) Create(c *gin.Context) {
	var req request.LibraryCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	library, err := libraryService.Create(&req)
	if err != nil {
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}

	response.OkWithData(library, c)
}

// Get 获取库详情（带统计信息）
func (l *LibraryApi) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	library, err := libraryService.GetLibraryInfo(uint(id))
	if err != nil {
		response.FailWithMessage("库不存在", c)
		return
	}

	response.OkWithData(library, c)
}

// Update 更新库
func (l *LibraryApi) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	var req request.LibraryCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	library, err := libraryService.Update(uint(id), &req)
	if err != nil {
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}

	response.OkWithData(library, c)
}

// Delete 删除库
func (l *LibraryApi) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	if err := libraryService.Delete(uint(id)); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.FailWithMessage("库不存在", c)
			return
		}
		response.FailWithMessage("删除失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// GetVersions 获取库的所有版本（用于上传时选择）
func (l *LibraryApi) GetVersions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	versions, err := libraryService.GetVersions(uint(id))
	if err != nil {
		response.FailWithMessage("获取版本列表失败: "+err.Error(), c)
		return
	}

	response.OkWithData(versions, c)
}

// CreateVersion 创建新版本
func (l *LibraryApi) CreateVersion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	var req request.VersionCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	if err := libraryService.CreateVersion(uint(id), req.Version); err != nil {
		response.FailWithMessage("创建版本失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("版本创建成功", c)
}

// DeleteVersion 删除版本及其所有文档
func (l *LibraryApi) DeleteVersion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	version := c.Param("version")
	if version == "" {
		response.FailWithMessage("版本不能为空", c)
		return
	}

	if err := libraryService.DeleteVersion(uint(id), version); err != nil {
		response.FailWithMessage("删除版本失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("版本删除成功", c)
}

// RefreshVersion 刷新版本（重新处理所有文档）
func (l *LibraryApi) RefreshVersion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	version := c.Param("version")
	if version == "" {
		response.FailWithMessage("版本不能为空", c)
		return
	}

	if err := libraryService.RefreshVersion(uint(id), version); err != nil {
		response.FailWithMessage("刷新版本失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("版本刷新已启动，请稍候", c)
}

// RefreshVersionSSE 刷新版本（SSE 实时推送处理状态）
func (l *LibraryApi) RefreshVersionSSE(c *gin.Context) {
	// 创建 SSE 写入器
	sse, ok := response.NewSSEWriter(c)
	if !ok {
		c.JSON(500, gin.H{"error": "SSE not supported"})
		return
	}

	// 解析参数
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		sse.SendError("无效的库ID")
		return
	}

	version := c.Param("version")
	if version == "" {
		sse.SendError("版本不能为空")
		return
	}

	// 创建状态通道
	statusChan := make(chan response.RefreshStatus, 20)

	// 启动后台处理
	go libraryService.RefreshVersionWithCallback(uint(id), version, statusChan)

	// 监听状态并推送 SSE
	for status := range statusChan {
		if status.Stage == "error" {
			sse.SendError(status.Message)
			return
		}
		sse.SendSuccess(status.Message, status)
	}
}
