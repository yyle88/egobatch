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
