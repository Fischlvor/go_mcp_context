package api

import (
	"errors"
	"strconv"

	"go-mcp-context/internal/model/response"
	"go-mcp-context/internal/service"

	"github.com/gin-gonic/gin"
)

type DocumentApi struct{}

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
