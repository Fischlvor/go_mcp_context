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

// List 获取库列表
func (l *LibraryApi) List(c *gin.Context) {
	var req request.LibraryList
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	result, err := libraryService.List(&req)
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

// Get 获取库详情
func (l *LibraryApi) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	library, err := libraryService.GetByID(uint(id))
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
