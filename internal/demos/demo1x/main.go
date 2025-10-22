package main

import (
	"context"
	"fmt"
	"time"

	"github.com/yyle88/egobatch/erxgroup"
)

// MyError is a custom error type with Code and Msg fields
// MyError 是具有 Code 和 Msg 字段的自定义错误类型
type MyError struct {
	Code string
	Msg  string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Msg)
}

func main() {
	ctx := context.Background()
	ego := erxgroup.NewGroup[*MyError](ctx)

	// Add task 1: takes 100ms to finish
	// 添加任务 1：需要 100ms 完成
	ego.Go(func(ctx context.Context) *MyError {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("Task 1 finished OK")
		return nil
	})

	// Add task 2: takes 50ms to finish
	// 添加任务 2：需要 50ms 完成
	ego.Go(func(ctx context.Context) *MyError {
		time.Sleep(50 * time.Millisecond)
		fmt.Println("Task 2 finished OK")
		return nil
	})

	// Add task 3: takes 80ms to finish
	// 添加任务 3：需要 80ms 完成
	ego.Go(func(ctx context.Context) *MyError {
		time.Sleep(80 * time.Millisecond)
		fmt.Println("Task 3 finished OK")
		return nil
	})

	// Wait until tasks finish and get the first issue
	// 等待所有任务完成并获取第一个问题（如果存在）
	if erx := ego.Wait(); erx != nil {
		fmt.Printf("Got issue: %s\n", erx.Error())
	} else {
		fmt.Println("Tasks finished OK")
	}
}
