package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/internal/service"
	"go-mcp-context/pkg/actlog"
	"go-mcp-context/pkg/global"
	"go-mcp-context/pkg/utils"

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

	// 记录活动日志（库创建成功后）
	userUUID := utils.GetUUID(c).String()
	actlog.Success(library.ID, actlog.EventLibCreate, fmt.Sprintf("创建库: %s", library.Name),
		actlog.WithActor(userUUID),
		actlog.WithTarget("library", fmt.Sprintf("%d", library.ID)),
	)

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

	// 记录活动日志
	actlog.Success(library.ID, actlog.EventLibUpdate, fmt.Sprintf("更新库: %s", library.Name))

	response.OkWithData(library, c)
}

// Delete 删除库
func (l *LibraryApi) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	// 删除前记录日志（删除后 library_id 仍然有效，因为是软删除）
	actlog.Info(uint(id), actlog.EventLibDelete, "开始删除库")

	if err := libraryService.Delete(uint(id)); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.FailWithMessage("库不存在", c)
			return
		}
		actlog.Error(uint(id), actlog.EventLibDelete, "删除库失败: "+err.Error())
		response.FailWithMessage("删除失败", c)
		return
	}

	actlog.Success(uint(id), actlog.EventLibDelete, "删除库成功")
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
		actlog.Error(uint(id), actlog.EventVerCreate, fmt.Sprintf("创建版本失败: %s - %s", req.Version, err.Error()),
			actlog.WithVersion(req.Version))
		response.FailWithMessage("创建版本失败: "+err.Error(), c)
		return
	}

	actlog.Success(uint(id), actlog.EventVerCreate, fmt.Sprintf("创建版本: %s", req.Version),
		actlog.WithVersion(req.Version))
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

	userUUID := utils.GetUUID(c).String()

	if err := libraryService.DeleteVersion(uint(id), version); err != nil {
		actlog.Error(uint(id), actlog.EventVerDelete, fmt.Sprintf("删除版本失败: %s - %s", version, err.Error()),
			actlog.WithActor(userUUID),
			actlog.WithVersion(version),
			actlog.WithTarget("version", version))
		response.FailWithMessage("删除版本失败: "+err.Error(), c)
		return
	}

	actlog.Success(uint(id), actlog.EventVerDelete, fmt.Sprintf("删除版本: %s", version),
		actlog.WithActor(userUUID),
		actlog.WithVersion(version),
		actlog.WithTarget("version", version))
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

	userUUID := utils.GetUUID(c).String()

	if err := libraryService.RefreshVersion(uint(id), version, userUUID); err != nil {
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

	// 获取操作者 ID
	actorID := utils.GetUUID(c).String()

	// 创建状态通道
	statusChan := make(chan response.RefreshStatus, 20)

	// 启动后台处理
	go libraryService.RefreshVersionWithCallback(uint(id), version, actorID, statusChan)

	// 监听状态并推送 SSE
	for status := range statusChan {
		if status.Stage == "error" {
			sse.SendError(status.Message)
			return
		}
		sse.SendSuccess(status.Message, status)
	}
}

// ==========================================
// GitHub 导入相关 API
// ==========================================

// ImportFromGitHubAsync 从 GitHub 导入文档（异步，立即返回）
// @Summary 从 GitHub 导入文档（异步）
// @Tags Library
// @Accept json
// @Produce json
// @Param id path int true "库 ID"
// @Param data body request.GitHubImportRequest true "导入参数"
// @Success 200 {object} response.Response
// @Router /libraries/github/import [post]
func (l *LibraryApi) ImportFromGitHub(c *gin.Context) {
	// 解析库 ID（从 query 参数获取）
	id, err := strconv.ParseUint(c.Query("id"), 10, 32)
	if err != nil || id == 0 {
		response.FailWithMessage("无效的库ID", c)
		return
	}

	// 解析请求参数
	var req request.GitHubImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	if req.Repo == "" {
		response.FailWithMessage("仓库名不能为空", c)
		return
	}

	// 检查版本是否已存在
	version := req.Version
	if version == "" {
		version = "latest"
	}
	library, err := libraryService.GetByID(uint(id))
	if err != nil {
		response.FailWithMessage("库不存在", c)
		return
	}
	for _, v := range library.Versions {
		if v == version {
			response.FailWithMessage(fmt.Sprintf("版本 %s 已存在", version), c)
			return
		}
	}

	// 同步写入"导入开始"日志（确保 API 返回前日志已入库）
	userUUID := utils.GetUUID(c).String()
	taskID := utils.GenerateTaskID()
	startLog := &dbmodel.ActivityLog{
		LibraryID:  uint(id),
		ActorID:    userUUID,
		TaskID:     taskID,
		Event:      actlog.EventGHImportStart,
		Status:     actlog.StatusInfo,
		Message:    fmt.Sprintf("开始导入: %s@%s", req.Repo, version),
		Version:    version,
		TargetType: "version",
		TargetID:   version,
	}
	global.DB.Create(startLog)

	// 启动后台导入（传递 taskID）
	req.TaskID = taskID
	githubService := service.NewGitHubImportService()
	go func() {
		githubService.ImportFromGitHub(context.Background(), uint(id), &req, userUUID, nil)
	}()

	response.OkWithMessage("GitHub 导入已启动，请通过活动日志查看进度", c)
}

// ImportFromGitHubSSE 从 GitHub 导入文档（SSE 实时推送进度）
func (l *LibraryApi) ImportFromGitHubSSE(c *gin.Context) {
	// 创建 SSE 写入器
	sse, ok := response.NewSSEWriter(c)
	if !ok {
		c.JSON(500, gin.H{"error": "SSE not supported"})
		return
	}

	// 解析库 ID（从 query 参数获取）
	id, err := strconv.ParseUint(c.Query("id"), 10, 32)
	if err != nil || id == 0 {
		sse.SendError("无效的库ID")
		return
	}

	// 解析请求参数
	var req request.GitHubImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sse.SendError("参数错误: " + err.Error())
		return
	}

	if req.Repo == "" {
		sse.SendError("仓库名不能为空")
		return
	}

	// 创建进度通道
	progressChan := make(chan response.GitHubImportProgress, 100)
	userUUID := utils.GetUUID(c).String()

	// 启动导入
	githubService := service.NewGitHubImportService()
	go func() {
		if err := githubService.ImportFromGitHub(c.Request.Context(), uint(id), &req, userUUID, progressChan); err != nil {
			// 错误已通过 progressChan 发送
		}
	}()

	// 监听进度并推送 SSE
	for progress := range progressChan {
		if progress.Stage == "failed" {
			sse.SendError(progress.Message)
			return
		}
		sse.SendSuccess(progress.Message, progress)
	}
}

// GetGitHubReleases 获取 GitHub 仓库的版本列表
func (l *LibraryApi) GetGitHubReleases(c *gin.Context) {
	repo := c.Query("repo")
	if repo == "" {
		response.FailWithMessage("仓库名不能为空", c)
		return
	}

	githubService := service.NewGitHubImportService()
	versions, err := githubService.GetMajorVersions(c.Request.Context(), repo, 20)
	if err != nil {
		response.FailWithMessage("获取版本失败: "+err.Error(), c)
		return
	}

	// 获取仓库信息
	repoInfo, err := githubService.GetRepoInfo(c.Request.Context(), repo)
	if err != nil {
		response.FailWithMessage("获取仓库信息失败: "+err.Error(), c)
		return
	}

	response.OkWithData(map[string]interface{}{
		"repo":           repo,
		"default_branch": repoInfo.DefaultBranch,
		"description":    repoInfo.Description,
		"versions":       versions,
	}, c)
}

// InitImportFromGitHub 从 GitHub URL 初始化导入（创建库 + 导入默认分支）
func (l *LibraryApi) InitImportFromGitHub(c *gin.Context) {
	var req request.GitHubInitImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 1. 初始化创建库（解析URL、验证连通性、检查重复、创建）
	library, defaultBranch, err := libraryService.InitFromGitHub(c.Request.Context(), req.GitHubURL)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 2. 记录活动日志
	userUUID := utils.GetUUID(c).String()
	taskID := utils.GenerateTaskID()

	actlog.Success(library.ID, actlog.EventLibCreate, fmt.Sprintf("从 GitHub 创建库: %s", library.Name),
		actlog.WithActor(userUUID),
		actlog.WithTarget("library", fmt.Sprintf("%d", library.ID)),
		actlog.WithVersion("latest"),
		actlog.WithTaskID(taskID),
	)

	// 3. 同步写入"导入开始"日志（确保 API 返回前日志已入库）
	startLog := &dbmodel.ActivityLog{
		LibraryID:  library.ID,
		ActorID:    userUUID,
		TaskID:     taskID,
		Event:      actlog.EventGHImportStart,
		Status:     actlog.StatusInfo,
		Message:    fmt.Sprintf("开始导入: %s@latest", library.SourceURL),
		Version:    "latest",
		TargetType: "version",
		TargetID:   "latest",
	}
	global.DB.Create(startLog)

	// 4. 异步导入（使用默认分支，版本名为 latest，传递 taskID）
	githubService := service.NewGitHubImportService()
	go func() {
		importReq := &request.GitHubImportRequest{
			Repo:    library.SourceURL,
			Branch:  defaultBranch,
			Version: "latest",
			TaskID:  taskID,
		}
		githubService.ImportFromGitHub(context.Background(), library.ID, importReq, userUUID, nil)
	}()

	response.OkWithData(response.GitHubInitImportResponse{
		LibraryID: library.ID,
		Version:   "latest",
	}, c)
}
