package egobatch

import (
	"context"

	"github.com/yyle88/egobatch/erxgroup"
	"github.com/yyle88/egobatch/internal/constraint"
	"github.com/yyle88/egobatch/internal/utils"
	"github.com/yyle88/must"
	"github.com/yyle88/must/mustnum"
)

// TaskBatch manages batch task execution with concurrent processing
// Supports glide mode enabling independent task execution and fail-fast mode
// Provides context error handling and result aggregation capabilities
//
// TaskBatch 管理批量任务的并发执行
// 支持平滑模式，可以独立执行任务或快速失败
// 提供上下文错误处理和结果聚合能力
type TaskBatch[A any, R any, E ErrorType] struct {
	Tasks Tasks[A, R, E]    // Task collection with arguments and results // 任务集合，包含参数和结果
	Glide bool              // Glide mode flag: false=fail-fast, true=independent tasks // 平滑模式标志：false=快速失败，true=独立任务
	waCtx func(err error) E // Context error conversion function // 上下文错误转换函数
}

// NewTaskBatch creates batch task engine with starting arguments
// Each argument becomes a task with zero-initialized result and error
// Default glide mode is false (fail-fast mode)
//
// NewTaskBatch 使用初始参数创建批量任务处理器
// 每个参数成为一个任务，结果和错误初始化为零值
// 默认平滑模式是 false（快速失败行为）
func NewTaskBatch[A any, R any, E ErrorType](args []A) *TaskBatch[A, R, E] {
	tasks := make([]*Task[A, R, E], 0, len(args))
	for idx := 0; idx < len(args); idx++ {
		tasks = append(tasks, &Task[A, R, E]{
			Arg: args[idx],
			Res: utils.Zero[R](),
			Erx: utils.Zero[E](),
		})
	}
	return &TaskBatch[A, R, E]{
		Tasks: tasks,
		Glide: false,
	}
}

// GetRun creates execution function at given index compatible with errgroup.Go
// Index must be valid (invoking code controls iteration count as basic contract)
// Returns wrapped function handling context cancellation and error propagation
//
// GetRun 在给定索引处创建与 errgroup.Go 兼容的执行函数
// 索引必须有效（调用者控制迭代次数作为基本约定）
// 返回处理上下文取消和错误传播的包装函数
func (t *TaskBatch[A, R, E]) GetRun(idx int, run func(ctx context.Context, arg A) (R, E)) func(ctx context.Context) E {
	mustnum.Less(idx, len(t.Tasks)) // Index bounds check - invoking code must not exceed task count // 索引边界检查 - 调用代码不能超过任务数量
	task := t.Tasks[idx]
	return func(ctx context.Context) E {
		if t.waCtx != nil && ctx.Err() != nil {
			erx := t.waCtx(ctx.Err()) // Convert context error - must return valid error, not fake zero // 转换上下文错误 - 必须返回有效错误，不能是伪造的零值
			must.False(constraint.Pass(erx))
			task.Erx = erx
			if t.Glide {
				return utils.Zero[E]() // Glide mode: record error but continue processing remaining tasks // 平滑模式：记录错误但继续处理剩余任务
			}
			return erx
		}
		res, erx := run(ctx, task.Arg) // Execute task - no panic allowed, invoking code must handle panic recovery // 执行任务 - 不允许 panic，调用代码必须处理 panic 恢复
		if !constraint.Pass(erx) {
			task.Erx = erx
			if t.Glide {
				return utils.Zero[E]() // Glide mode: record error without canceling context, allowing other tasks to proceed // 平滑模式：记录错误但不取消上下文，允许其他任务继续
			}
			return erx
		}
		task.Res = res
		return utils.Zero[E]()
	}
}

// EgoRun demonstrates GetRun usage with inversion-of-control pattern
// When task logic is complex and scheduling logic is simple, pass scheduling engine as argument
// Auto schedules tasks into the provided errgroup
//
// EgoRun 演示 GetRun 使用方式，采用控制反转模式
// 当任务逻辑较重而调度逻辑较轻时，将调度器作为参数传入
// 自动将所有任务调度到提供的 errgroup 中
func (t *TaskBatch[A, R, E]) EgoRun(ego *erxgroup.Group[E], run func(ctx context.Context, arg A) (R, E)) {
	for idx := 0; idx < len(t.Tasks); idx++ {
		ego.Go(t.GetRun(idx, run))
	}
}

// SetGlide configures glide mode
// When true: tasks execute in independent mode, errors recorded without stopping others
// When false: first error stops batch execution (fail-fast)
//
// SetGlide 配置平滑模式行为
// 当为 true：任务独立执行，错误被记录但不停止其他任务
// 当为 false：第一个错误停止批量执行（快速失败）
func (t *TaskBatch[A, R, E]) SetGlide(glide bool) {
	t.Glide = glide
}

// SetWaCtx configures context error conversion function
// Converts context.Context errors into custom error type E
// Gets invoked when context cancellation happens
//
// SetWaCtx 配置上下文错误转换函数
// 将 context.Context 错误转换为自定义错误类型 E
// 在上下文取消或超时发生时调用
func (t *TaskBatch[A, R, E]) SetWaCtx(waCtx func(err error) E) {
	t.waCtx = waCtx
}
