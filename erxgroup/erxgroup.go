// Package erxgroup provides generic wrapper around errgroup with type-safe custom error handling
// Enables using custom error types instead of standard error interface while maintaining errgroup semantics
//
// Package erxgroup 提供 errgroup 的泛型包装，支持类型安全的自定义错误处理
// 允许使用自定义错误类型而非标准 error 接口，同时保持 errgroup 语义
package erxgroup

import (
	"context"
	"errors"

	"github.com/yyle88/egobatch/internal/constraint"
	"github.com/yyle88/egobatch/internal/utils"
	"github.com/yyle88/must"
	"golang.org/x/sync/errgroup"
)

// ErrorType is an alias to the constraint defined in the internal package
// ErrorType 是内部包中定义的约束的别名
type ErrorType = constraint.ErrorType

// Group wraps errgroup.Group with type-safe custom error handling
// Provides generic error type E instead of standard error interface
// Maintains context cancellation and goroutine synchronization semantics
//
// Group 使用类型安全的自定义错误处理包装 errgroup.Group
// 提供泛型错误类型 E 而不是标准 error 接口
// 保持上下文取消和协程同步语义
type Group[E ErrorType] struct {
	ego *errgroup.Group // Underlying errgroup instance // 底层 errgroup 实例
	ctx context.Context // Shared context with cancellation // 共享的可取消上下文
}

// NewGroup creates generic errgroup with custom error type
// Context cancels when first error occurs or parent context cancels
//
// NewGroup 创建带有自定义错误类型的泛型 errgroup
// 当第一个错误发生或父上下文取消时，上下文会被取消
func NewGroup[E ErrorType](ctx context.Context) *Group[E] {
	ego, ctx := errgroup.WithContext(ctx)
	return &Group[E]{
		ego: ego,
		ctx: ctx,
	}
}

// Wait blocks awaiting goroutine completion and returns first error
// Uses errors.As to convert standard error back to custom type E
// Returns zero value when execution succeeds
//
// Wait 阻塞直到所有协程完成并返回第一个错误
// 使用 errors.As 将标准 error 转换回自定义类型 E
// 当所有协程成功时返回零值
func (G *Group[E]) Wait() E {
	if err := G.ego.Wait(); err != nil {
		var erx E
		must.True(errors.As(err, &erx))
		return erx
	}
	return utils.Zero[E]()
}

// Go starts goroutine within the group
// Converts custom error E to standard error when task fails
// Task receives shared cancellable context
//
// Go 在组内启动协程
// 当任务失败时将自定义错误 E 转换为标准 error
// 任务接收共享的可取消上下文
func (G *Group[E]) Go(run func(ctx context.Context) E) {
	G.ego.Go(func() error {
		if erx := run(G.ctx); !constraint.Pass(erx) {
			return erx
		}
		return nil
	})
}

// TryGo attempts to start goroutine within the group
// Returns false if goroutine limit reached, true if started
// Same error handling as Go method
//
// TryGo 尝试在组内启动协程
// 如果达到协程限制则返回 false，如果启动则返回 true
// 与 Go 方法相同的错误处理
func (G *Group[E]) TryGo(run func(ctx context.Context) E) bool {
	return G.ego.TryGo(func() error {
		if erx := run(G.ctx); !constraint.Pass(erx) {
			return erx
		}
		return nil
	})
}

// SetLimit restricts concurrent goroutines count
// Must be invoked before first Go and TryGo invocation
//
// SetLimit 限制并发协程数量
// 必须在第一次 Go 或 TryGo 调用之前调用
func (G *Group[E]) SetLimit(n int) {
	G.ego.SetLimit(n)
}
