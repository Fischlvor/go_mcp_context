package service

import "errors"

// 服务层通用错误
var (
	ErrNotFound      = errors.New("资源不存在")
	ErrAlreadyExists = errors.New("资源已存在")
	ErrInvalidParams = errors.New("参数无效")
	ErrUnauthorized  = errors.New("未授权")
	ErrForbidden     = errors.New("禁止访问")
)
