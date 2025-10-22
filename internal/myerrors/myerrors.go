package myerrors

import (
	"errors"
	"fmt"
)

// Error is a simple custom error type that is comparable
// Error 是一个简单的自定义错误类型，是可比较的
type Error struct {
	code string
	desc string
}

// Error implements the error interface
// Error 实现 error 接口
func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s", e.code, e.desc)
}

// Code returns the error code
// Code 返回错误码
func (e *Error) Code() string {
	return e.code
}

// New creates a new error with code and message
// New 创建带有错误码和消息的新错误
func New(code string, format string, args ...interface{}) *Error {
	return &Error{
		code: code,
		desc: fmt.Sprintf(format, args...),
	}
}

// ErrorServiceError creates a service error
// ErrorServiceError 创建服务错误
func ErrorServiceError(format string, args ...interface{}) *Error {
	return New("SERVICE_ERROR", format, args...)
}

// ErrorWrongContext creates a context error
// ErrorWrongContext 创建上下文错误
func ErrorWrongContext(format string, args ...interface{}) *Error {
	return New("CONTEXT_ERROR", format, args...)
}

// IsServiceError checks if the error is a service error
// IsServiceError 检查错误是否是服务错误
func IsServiceError(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.code == "SERVICE_ERROR"
	}
	return false
}

// IsWrongContext checks if the error is a context error
// IsWrongContext 检查错误是否是上下文错误
func IsWrongContext(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.code == "CONTEXT_ERROR"
	}
	return false
}
