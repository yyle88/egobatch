package main

import (
	"context"
	"fmt"
	"time"

	"github.com/yyle88/egobatch"
	"github.com/yyle88/egobatch/erxgroup"
	"github.com/yyle88/must"
)

// MyError is a custom error type with error code
// MyError 是带有错误代码的自定义错误类型
type MyError struct {
	Code string
}

func (e *MyError) Error() string {
	return e.Code
}

func main() {
	// Create context with 150ms timeout
	// 创建带 150ms 超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	// Create batch with task arguments
	// 创建批量任务参数
	args := []int{1, 2, 3, 4, 5}
	batch := egobatch.NewTaskBatch[int, string, *MyError](args)

	// Use glide mode to see which tasks finish vs timeout
	// 使用平滑模式观察哪些任务完成、哪些超时
	batch.SetGlide(true)

	// Convert context issues to custom error type
	// 将上下文问题转换为自定义错误类型
	batch.SetWaCtx(func(err error) *MyError {
		return &MyError{Code: "TIMEOUT"}
	})

	ego := erxgroup.NewGroup[*MyError](ctx)

	batch.EgoRun(ego, func(ctx context.Context, num int) (string, *MyError) {
		// Each task needs different time: 50ms, 100ms, 150ms, 200ms, 250ms
		// 每个任务需要不同时间：50ms、100ms、150ms、200ms、250ms
		taskTime := time.Duration(num*50) * time.Millisecond

		timer := time.NewTimer(taskTime)
		defer timer.Stop()

		select {
		case <-timer.C:
			// Task finishes within timeout
			// 任务在超时前完成
			fmt.Printf("Task %d finished (%dms)\n", num, num*50)
			return fmt.Sprintf("task-%d", num), nil
		case <-ctx.Done():
			// Task cancelled due to timeout
			// 任务因超时而取消
			fmt.Printf("Task %d cancelled (%dms needed)\n", num, num*50)
			return "", &MyError{Code: "CANCELLED"}
		}
	})

	// In glide mode, ego.Wait() returns nil because errors are captured in tasks
	// 在平滑模式下，ego.Wait() 返回 nil 因为错误已被捕获在任务中
	must.Null(ego.Wait())

	// Show task outcomes
	// 显示任务结果
	okTasks := batch.Tasks.OkTasks()
	waTasks := batch.Tasks.WaTasks()

	fmt.Printf("\nSuccess: %d, Timeout: %d\n", len(okTasks), len(waTasks))

	// Show finished tasks
	// 显示完成的任务
	for _, task := range okTasks {
		fmt.Printf("Arg: %d -> Outcome: %s\n", task.Arg, task.Res)
	}

	// Show timed-out tasks
	// 显示超时的任务
	for _, task := range waTasks {
		fmt.Printf("Arg: %d -> Issue: %s\n", task.Arg, task.Erx.Error())
	}
}
