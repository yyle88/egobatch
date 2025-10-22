// Package example3 demonstrates multi-step pipeline processing patterns
// Shows cascading task execution with nested batch operations
//
// 包 example3 演示多步骤流水线处理模式
// 展示级联任务执行以及嵌套批处理操作
package example3

import (
	"github.com/yyle88/egobatch"
	"github.com/yyle88/egobatch/internal/myerrors"
)

// Step1Param represents input arguments on step 1
// 步骤1的输入参数
type Step1Param struct {
	NumA int // Input value A // 输入值A
}

// Step1Result represents step 1 result with nested step 2 outputs
// 步骤1结果，包含嵌套的步骤2输出
type Step1Result struct {
	ResA         string                                                              // Result value A // 结果值A
	Step2Outputs egobatch.TaskOutputList[*Step2Param, *Step2Result, *myerrors.Error] // Nested step 2 outputs // 嵌套的步骤2输出
}

// Step2Param represents input arguments on step 2
// 步骤2的输入参数
type Step2Param struct {
	NumB int // Input value B // 输入值B
}

// Step2Result represents step 2 result with nested step 3 outputs
// 步骤2结果，包含嵌套的步骤3输出
type Step2Result struct {
	ResB         string                                                              // Result value B // 结果值B
	Step3Outputs egobatch.TaskOutputList[*Step3Param, *Step3Result, *myerrors.Error] // Nested step 3 outputs // 嵌套的步骤3输出
}

// Step3Param represents input arguments on step 3
// 步骤3的输入参数
type Step3Param struct {
	NumC int // Input value C // 输入值C
}

// Step3Result represents step 3 final result
// 步骤3的最终结果
type Step3Result struct {
	ResC string // Result value C // 结果值C
}

// NewStep1Params creates a collection of step 1 arguments
// Values assigned with sequential index
//
// NewStep1Params 创建步骤1的参数集合
// 数值按索引赋值
func NewStep1Params(paramCount int) []*Step1Param {
	var params = make([]*Step1Param, 0, paramCount)
	for idx := 0; idx < paramCount; idx++ {
		params = append(params, &Step1Param{NumA: idx})
	}
	return params
}

// NewStep2Params creates a collection of step 2 arguments
// Values assigned with sequential index
//
// NewStep2Params 创建步骤2的参数集合
// 数值按索引赋值
func NewStep2Params(paramCount int) []*Step2Param {
	var params = make([]*Step2Param, 0, paramCount)
	for idx := 0; idx < paramCount; idx++ {
		params = append(params, &Step2Param{NumB: idx})
	}
	return params
}

// NewStep3Params creates a collection of step 3 arguments
// Values assigned with sequential index
//
// NewStep3Params 创建步骤3的参数集合
// 数值按索引赋值
func NewStep3Params(paramCount int) []*Step3Param {
	var params = make([]*Step3Param, 0, paramCount)
	for idx := 0; idx < paramCount; idx++ {
		params = append(params, &Step3Param{NumC: idx})
	}
	return params
}
