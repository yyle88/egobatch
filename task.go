// Package egobatch: Generic batch task processing with type-safe error handling
// Provides Task and TaskBatch types enabling concurrent batch operations with custom error types
// Supports result filtering (OkTasks/WaTasks) and flexible result transformation via Flatten
//
// egobatch: 具有类型安全错误处理的泛型批量任务处理
// 提供 Task 和 TaskBatch 类型，支持使用自定义错误类型的并发批量操作
// 支持结果过滤（OkTasks/WaTasks）和通过 Flatten 进行灵活的结果转换
package egobatch

import (
	"github.com/yyle88/egobatch/internal/constraint"
)

// ErrorType is an alias to the constraint defined in the internal package
// Enables type-safe error handling with comparable custom error types
//
// ErrorType 是内部包中定义的约束的别名
// 支持使用可比较的自定义错误类型进行类型安全的错误处理
type ErrorType = constraint.ErrorType

// Task represents a single task with argument, result, and error
// Generic type supporting any argument type A, result type R, and error type E
//
// Task 代表单个任务，包含参数、结果和错误
// 泛型类型支持任意参数类型 A、结果类型 R 和错误类型 E
type Task[A any, R any, E ErrorType] struct {
	Arg A // Task input argument // 任务输入参数
	Res R // Task result value // 任务结果值
	Erx E // Task error (nil when success) // 任务错误（成功时为 nil）
}

// Tasks is a slice of Task pointers supporting batch operations
// Provides filtering and transformation methods on task collections
//
// Tasks 是 Task 指针切片，支持批量操作
// 提供任务集合的过滤和转换方法
type Tasks[A any, R any, E ErrorType] []*Task[A, R, E]

// OkTasks filters and returns tasks that completed with success
// Returns subset of tasks on success
//
// OkTasks 过滤并返回成功完成的任务
// 返回成功的任务子集
func (tasks Tasks[A, R, E]) OkTasks() Tasks[A, R, E] {
	var okTasks Tasks[A, R, E]
	for _, task := range tasks {
		if constraint.Pass(task.Erx) {
			okTasks = append(okTasks, task)
		}
	}
	return okTasks
}

// WaTasks filters and returns tasks that failed with errors
// Returns subset of tasks when error occurs
//
// WaTasks 过滤并返回失败的任务
// 返回出错的任务子集
func (tasks Tasks[A, R, E]) WaTasks() Tasks[A, R, E] {
	var waTasks Tasks[A, R, E]
	for _, task := range tasks {
		if !constraint.Pass(task.Erx) {
			waTasks = append(waTasks, task)
		}
	}
	return waTasks
}

// Flatten transforms task results into flat slice with error handling
// Uses newWaFunc to convert failed tasks into result type R
// Returns slice of results mixing success cases and transformed errors
//
// Flatten 将任务结果转换成扁平切片并处理错误
// 使用 newWaFunc 将失败的任务转换为结果类型 R
// 返回混合成功结果和转换后错误的结果切片
func (tasks Tasks[A, R, E]) Flatten(newWaFunc func(arg A, erx E) R) []R {
	var results = make([]R, 0, len(tasks))
	for _, task := range tasks {
		if !constraint.Pass(task.Erx) {
			results = append(results, newWaFunc(task.Arg, task.Erx))
		} else {
			results = append(results, task.Res)
		}
	}
	return results
}
