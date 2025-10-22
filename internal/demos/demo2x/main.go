package main

import (
	"context"
	"fmt"

	"github.com/yyle88/egobatch"
	"github.com/yyle88/egobatch/erxgroup"
	"github.com/yyle88/must"
)

// MyError is a simple custom error type with error code
// MyError 是带有错误代码的简单自定义错误类型
type MyError struct {
	Code string
}

func (e *MyError) Error() string {
	return e.Code
}

func main() {
	// Create batch with arguments
	// 使用参数创建批量任务
	args := []int{1, 2, 3, 4, 5}
	batch := egobatch.NewTaskBatch[int, string, *MyError](args)

	// Configure glide mode - keep going even when issues happen
	// 配置平滑模式 - 即使出现问题也继续处理
	batch.SetGlide(true)

	// Execute batch tasks
	// 执行批量任务
	ctx := context.Background()
	ego := erxgroup.NewGroup[*MyError](ctx)

	batch.EgoRun(ego, func(ctx context.Context, num int) (string, *MyError) {
		if num%2 == 0 {
			// Even numbers finish OK
			// 偶数处理完成
			return fmt.Sprintf("even-%d", num), nil
		}
		// Odd numbers have issues
		// 奇数出现问题
		return "", &MyError{Code: "ODD_NUMBER"}
	})

	// In glide mode, ego.Wait() returns nil because errors are captured in tasks
	// 在平滑模式下，ego.Wait() 返回 nil 因为错误已被捕获在任务中
	must.Null(ego.Wait())

	// Get and handle task outcomes
	// 获取和处理任务结果
	okTasks := batch.Tasks.OkTasks()
	waTasks := batch.Tasks.WaTasks()

	fmt.Printf("Success: %d, Failed: %d\n", len(okTasks), len(waTasks))

	// Show OK outcomes
	// 显示成功结果
	for _, task := range okTasks {
		fmt.Printf("Arg: %d -> Outcome: %s\n", task.Arg, task.Res)
	}

	// Show bad outcomes
	// 显示失败结果
	for _, task := range waTasks {
		fmt.Printf("Arg: %d -> Issue: %s\n", task.Arg, task.Erx.Error())
	}
}
