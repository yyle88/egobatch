package egobatch

import (
	"github.com/yyle88/egobatch/internal/constraint"
	"github.com/yyle88/egobatch/internal/utils"
)

// TaskOutput represents single task execution outcome with argument, result, and error
// Like Task but designed as simple data container without batch processing logic
// Used when returning task results without needing complete Task batch capabilities
//
// TaskOutput 代表单个任务执行结果，包含参数、结果和错误
// 类似 Task 但设计为简单数据传输对象，不包含批处理逻辑
// 用于返回任务结果而不需要完整的 Task 批处理能力
type TaskOutput[ARG any, RES any, E ErrorType] struct {
	Arg ARG // Task input argument // 任务输入参数
	Res RES // Task result value // 任务结果值
	Erx E   // Task error (nil when success) // 任务错误（成功时为 nil）
}

// NewOkTaskOutput creates success task output with result
// Error field initialized to zero value indicating success
//
// NewOkTaskOutput 创建成功的任务输出，包含结果
// 错误字段初始化为零值，表示成功
func NewOkTaskOutput[ARG any, RES any, E ErrorType](arg ARG, res RES) *TaskOutput[ARG, RES, E] {
	return &TaskOutput[ARG, RES, E]{
		Arg: arg,
		Res: res,
		Erx: utils.Zero[E](),
	}
}

// NewWaTaskOutput creates failed task output with error
// Result field initialized to zero value as task failed
//
// NewWaTaskOutput 创建失败的任务输出，包含错误
// 结果字段初始化为零值，因为任务失败
func NewWaTaskOutput[ARG any, RES any, E ErrorType](arg ARG, erx E) *TaskOutput[ARG, RES, E] {
	return &TaskOutput[ARG, RES, E]{
		Arg: arg,
		Res: utils.Zero[RES](),
		Erx: erx,
	}
}

// TaskOutputList is a collection of task outputs supporting filtering and aggregation
// Provides methods to separate success and failed results
// Enables result extraction and error collection patterns
//
// TaskOutputList 是任务输出的集合，支持过滤和聚合
// 提供分离成功和失败结果的方法
// 支持结果提取和错误收集模式
type TaskOutputList[ARG any, RES any, E ErrorType] []*TaskOutput[ARG, RES, E]

// OkList filters and returns outputs that completed with success
// Returns subset on success
//
// OkList 过滤并返回成功完成的输出
// 返回成功的子集
func (rs TaskOutputList[ARG, RES, E]) OkList() TaskOutputList[ARG, RES, E] {
	var results TaskOutputList[ARG, RES, E]
	for _, one := range rs {
		if constraint.Pass(one.Erx) {
			results = append(results, one)
		}
	}
	return results
}

// WaList filters and returns outputs that failed with errors
// Returns subset when error occurs
//
// WaList 过滤并返回失败的输出
// 返回出错的子集
func (rs TaskOutputList[ARG, RES, E]) WaList() TaskOutputList[ARG, RES, E] {
	var results TaskOutputList[ARG, RES, E]
	for _, one := range rs {
		if !constraint.Pass(one.Erx) {
			results = append(results, one)
		}
	}
	return results
}

// OkCount counts success task outputs
// Returns count of outputs on success
//
// OkCount 统计成功的任务输出
// 返回成功的输出数量
func (rs TaskOutputList[ARG, RES, E]) OkCount() int {
	var cnt int
	for _, one := range rs {
		if constraint.Pass(one.Erx) {
			cnt++
		}
	}
	return cnt
}

// WaCount counts failed task outputs
// Returns count of outputs when error occurs
//
// WaCount 统计失败的任务输出
// 返回出错的输出数量
func (rs TaskOutputList[ARG, RES, E]) WaCount() int {
	var cnt int
	for _, one := range rs {
		if !constraint.Pass(one.Erx) {
			cnt++
		}
	}
	return cnt
}

// OkResults extracts result values from success outputs
// Returns slice containing just results from outputs without errors
//
// OkResults 从成功的输出中提取结果值
// 返回仅包含无错误输出的结果切片
func (rs TaskOutputList[ARG, RES, E]) OkResults() []RES {
	var results []RES
	for _, one := range rs {
		if constraint.Pass(one.Erx) {
			results = append(results, one.Res)
		}
	}
	return results
}

// WaReasons extracts error values from failed outputs
// Returns slice containing just errors from outputs with failures
//
// WaReasons 从失败的输出中提取错误值
// 返回仅包含失败输出的错误切片
func (rs TaskOutputList[ARG, RES, E]) WaReasons() []E {
	var reasons []E
	for _, one := range rs {
		if !constraint.Pass(one.Erx) {
			reasons = append(reasons, one.Erx)
		}
	}
	return reasons
}
