// Package utils provides common utilities
// 包 utils 提供通用工具
package utils

// Zero returns zero value of type T using named return pattern
// Works with any type including interfaces, structs, primitives
//
// Zero 使用命名返回模式返回类型 T 的零值
// 适用于任何类型，包括接口、结构体、基本类型
func Zero[T any]() (x T) {
	return x
}

// Same checks if two values of comparable type T are same using direct comparison
// Requires T to be comparable to allow == check
//
// Same 使用直接比较检查可比较类型 T 的两个值是否相同
// 需要 T 可比较以允许使用 == 检查
func Same[T comparable](a, b T) bool {
	return a == b
}
