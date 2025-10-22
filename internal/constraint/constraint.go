package constraint

import "github.com/yyle88/egobatch/internal/utils"

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
// Uses direct comparison which is safe because ErrorType requires comparable
//
// Pass 检查错误是否为 nil（零值）
// 使用直接比较是安全的，因为 ErrorType 要求可比较
func Pass[E ErrorType](erx E) bool {
	// IMPORTANT: Do NOT use errors.Is here as it causes nil pointer panic
	// When E is pointer type, calling errors.Is(nonNilErr, nil) triggers Is(nil) method
	// Most error types (e.g. Kratos Error, and others) Is method don't check if self is nil, causing panic
	// Direct == comparison is safe because ErrorType requires comparable
	//
	// 重要：禁止使用 errors.Is 会导致 nil 指针 panic
	// 当 E 是指针类型时，errors.Is(非nil错误, nil) 会触发 Is(nil) 方法调用
	// 绝大多数错误类型（如 Kratos Error 等）的 Is 方法都不会再检查自身为 nil 因此会导致 panic
	// 直接 == 比较安全因为 ErrorType 要求可比较
	return utils.Same(erx, utils.Zero[E]())
}
