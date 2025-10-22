package constraint

import "errors"

// ErrorType defines a comparable error type constraint
// This allows using any pointer error type as the generic error parameter
//
// ErrorType 定义可比较的错误类型约束
// 允许使用任何指针错误类型作为泛型错误参数
type ErrorType interface {
	error
	comparable
}

// Pass checks if error is nil (zero value)
// Uses errors.Is to ensure correct nil checking with error wrapping support
//
// Pass 检查错误是否为 nil（零值）
// 使用 errors.Is 确保正确的 nil 检查，支持错误包装
func Pass[E ErrorType](erx E) bool {
	var zero E
	return errors.Is(erx, zero)
}
