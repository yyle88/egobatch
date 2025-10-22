[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yyle88/egobatch/release.yml?branch=main&label=BUILD)](https://github.com/yyle88/egobatch/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yyle88/egobatch)](https://pkg.go.dev/github.com/yyle88/egobatch)
[![Coverage Status](https://img.shields.io/coveralls/github/yyle88/egobatch/main.svg)](https://coveralls.io/github/yyle88/egobatch?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25%2B-lightgrey.svg)](https://github.com/yyle88/egobatch)
[![GitHub Release](https://img.shields.io/github/release/yyle88/egobatch.svg)](https://github.com/yyle88/egobatch/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yyle88/egobatch)](https://goreportcard.com/report/github.com/yyle88/egobatch)

# egobatch

具有类型安全自定义错误处理的泛型批量任务处理包。

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->
## 英文文档

[ENGLISH README](README.md)
<!-- TEMPLATE (ZH) END: LANGUAGE NAVIGATION -->

## 核心特性

🎯 **类型安全错误处理**: 支持可比较约束的泛型错误类型
⚡ **批量任务处理**: 使用 TaskBatch 抽象进行并发执行
🔄 **灵活执行模式**: 平滑模式（独立任务）和快速失败模式
🌍 **上下文取消支持**: 完整的上下文传播和超时支持
📋 **结果聚合**: 使用 OkTasks/WaTasks 方法过滤成功/失败任务

## 安装

```bash
go get github.com/yyle88/egobatch
```

## 快速开始

### 基础 errgroup 与自定义错误类型

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/yyle88/egobatch/erxgroup"
)

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

	// 添加任务 1：需要 100ms 完成
	ego.Go(func(ctx context.Context) *MyError {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("Task 1 finished OK")
		return nil
	})

	// 添加任务 2：需要 50ms 完成
	ego.Go(func(ctx context.Context) *MyError {
		time.Sleep(50 * time.Millisecond)
		fmt.Println("Task 2 finished OK")
		return nil
	})

	// 添加任务 3：需要 80ms 完成
	ego.Go(func(ctx context.Context) *MyError {
		time.Sleep(80 * time.Millisecond)
		fmt.Println("Task 3 finished OK")
		return nil
	})

	// 等待所有任务完成并获取第一个问题（如果存在）
	if erx := ego.Wait(); erx != nil {
		fmt.Printf("Got issue: %s\n", erx.Error())
	} else {
		fmt.Println("Tasks finished OK")
	}
}
```

⬆️ **源码:** [源码](internal/demos/demo1x/main.go)

### 批量任务处理

```go
package main

import (
	"context"
	"fmt"

	"github.com/yyle88/egobatch"
	"github.com/yyle88/egobatch/erxgroup"
	"github.com/yyle88/must"
)

// MyError 是带有错误代码的简单自定义错误类型
type MyError struct {
	Code string
}

func (e *MyError) Error() string {
	return e.Code
}

func main() {
	// 使用参数创建批量任务
	args := []int{1, 2, 3, 4, 5}
	batch := egobatch.NewTaskBatch[int, string, *MyError](args)

	// 配置平滑模式 - 即使出现问题也继续处理
	batch.SetGlide(true)

	// 执行批量任务
	ctx := context.Background()
	ego := erxgroup.NewGroup[*MyError](ctx)

	batch.EgoRun(ego, func(ctx context.Context, num int) (string, *MyError) {
		if num%2 == 0 {
			// 偶数处理完成
			return fmt.Sprintf("even-%d", num), nil
		}
		// 奇数出现问题
		return "", &MyError{Code: "ODD_NUMBER"}
	})

	// 在平滑模式下，ego.Wait() 返回 nil 因为错误已被捕获在任务中
	must.Null(ego.Wait())

	// 获取和处理任务结果
	okTasks := batch.Tasks.OkTasks()
	waTasks := batch.Tasks.WaTasks()

	fmt.Printf("Success: %d, Failed: %d\n", len(okTasks), len(waTasks))

	// 显示成功结果
	for _, task := range okTasks {
		fmt.Printf("Arg: %d -> Outcome: %s\n", task.Arg, task.Res)
	}

	// 显示失败结果
	for _, task := range waTasks {
		fmt.Printf("Arg: %d -> Issue: %s\n", task.Arg, task.Erx.Error())
	}
}
```

⬆️ **源码:** [源码](internal/demos/demo2x/main.go)

### 上下文超时处理

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/yyle88/egobatch"
	"github.com/yyle88/egobatch/erxgroup"
	"github.com/yyle88/must"
)

// MyError 是带有错误代码的自定义错误类型
type MyError struct {
	Code string
}

func (e *MyError) Error() string {
	return e.Code
}

func main() {
	// 创建带 150ms 超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	// 创建批量任务参数
	args := []int{1, 2, 3, 4, 5}
	batch := egobatch.NewTaskBatch[int, string, *MyError](args)

	// 使用平滑模式观察哪些任务完成、哪些超时
	batch.SetGlide(true)

	// 将上下文问题转换为自定义错误类型
	batch.SetWaCtx(func(err error) *MyError {
		return &MyError{Code: "TIMEOUT"}
	})

	ego := erxgroup.NewGroup[*MyError](ctx)

	batch.EgoRun(ego, func(ctx context.Context, num int) (string, *MyError) {
		// 每个任务需要不同时间：50ms、100ms、150ms、200ms、250ms
		taskTime := time.Duration(num*50) * time.Millisecond

		timer := time.NewTimer(taskTime)
		defer timer.Stop()

		select {
		case <-timer.C:
			// 任务在超时前完成
			fmt.Printf("Task %d finished (%dms)\n", num, num*50)
			return fmt.Sprintf("task-%d", num), nil
		case <-ctx.Done():
			// 任务因超时而取消
			fmt.Printf("Task %d cancelled (%dms needed)\n", num, num*50)
			return "", &MyError{Code: "CANCELLED"}
		}
	})

	// 在平滑模式下，ego.Wait() 返回 nil 因为错误已被捕获在任务中
	must.Null(ego.Wait())

	// 显示任务结果
	okTasks := batch.Tasks.OkTasks()
	waTasks := batch.Tasks.WaTasks()

	fmt.Printf("\nSuccess: %d, Timeout: %d\n", len(okTasks), len(waTasks))

	// 显示完成的任务
	for _, task := range okTasks {
		fmt.Printf("Arg: %d -> Outcome: %s\n", task.Arg, task.Res)
	}

	// 显示超时的任务
	for _, task := range waTasks {
		fmt.Printf("Arg: %d -> Issue: %s\n", task.Arg, task.Erx.Error())
	}
}
```

⬆️ **源码:** [源码](internal/demos/demo3x/main.go)

### 快速失败模式

```go
batch := egobatch.NewTaskBatch[int, string, *MyError](args)
// 默认是快速失败模式 (Glide: false)

ego := erxgroup.NewGroup[*MyError](ctx)
batch.EgoRun(ego, taskFunc)

if erx := ego.Wait(); erx != nil {
    // 第一个错误停止执行
    fmt.Printf("遇到错误停止: %s\n", erx.Error())
}
```

### 任务结果转换

```go
tasks := batch.Tasks

// 使用错误处理进行扁平化
results := tasks.Flatten(func(arg int, err *MyError) string {
    return fmt.Sprintf("错误-%d: %s", arg, err.Code)
})

// 混合成功结果和转换后的错误
for _, result := range results {
    fmt.Println(result)
}
```

## 核心组件

### erxgroup.Group[E ErrorType]

`errgroup.Group` 的泛型包装，具有类型安全的自定义错误：

- `NewGroup[E](ctx)`: 使用自定义错误类型创建新组
- `Go(func(ctx) E)`: 添加返回自定义错误的任务
- `TryGo(func(ctx) E)`: 添加带限制检查的任务
- `Wait() E`: 等待并获取第一个类型化错误
- `SetLimit(n)`: 限制并发任务数量

### TaskBatch[A, R, E]

批量任务执行与并发处理：

- `NewTaskBatch[A, R, E](args)`: 从参数创建批量任务
- `SetGlide(bool)`: 配置执行模式
- `SetWaCtx(func(error) E)`: 处理上下文错误
- `GetRun(idx, func)`: 获取任务执行函数
- `EgoRun(ego, func)`: 使用 errgroup 运行批量任务

### Tasks[A, R, E]

任务集合与过滤方法：

- `OkTasks()`: 获取成功完成的任务
- `WaTasks()`: 获取失败的任务
- `Flatten(func)`: 使用错误处理转换结果

## 高级用法

### 上下文超时处理

```go
batch := egobatch.NewTaskBatch[int, string, *MyError](args)
batch.SetGlide(true)

// 将上下文错误转换为自定义类型
batch.SetWaCtx(func(err error) *MyError {
    return &MyError{Code: "CONTEXT_ERROR", Msg: err.Error()}
})

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

ego := erxgroup.NewGroup[*MyError](ctx)
batch.EgoRun(ego, taskFunc)
ego.Wait()

// 超时的任务会记录上下文错误
for _, task := range batch.Tasks.WaTasks() {
    fmt.Printf("任务 %d 错误: %s\n", task.Arg, task.Erx.Error())
}
```

### TaskOutput 模式

```go
import "github.com/yyle88/egobatch"

// 创建任务输出
outputs := egobatch.TaskOutputList[int, string, *MyError]{
    egobatch.NewOkTaskOutput(1, "成功-1"),
    egobatch.NewWaTaskOutput(2, &MyError{Code: "FAIL"}),
    egobatch.NewOkTaskOutput(3, "成功-3"),
}

// 过滤和聚合
okList := outputs.OkList()
okCount := outputs.OkCount()
results := outputs.OkResults()
errors := outputs.WaReasons()

fmt.Printf("成功数量: %d\n", okCount)
fmt.Printf("结果: %v\n", results)
```

## 设计模式

### ErrorType 约束

自定义错误类型必须满足 `ErrorType` 约束：

```go
type ErrorType interface {
    error
    comparable
}
```

这使得：
- 使用 `errors.Is` 进行类型安全错误检查
- 使用 `constraint.Pass(erx)` 进行零值 nil 检测
- 使用 `errors.As` 进行自定义错误转换

### 平滑模式 vs 快速失败

**平滑模式 (Glide: true)**:
- 任务以独立模式执行
- 记录错误但不停止其他任务
- 上下文取消影响剩余任务
- 适合独立操作

**快速失败模式 (Glide: false)**:
- 第一个错误停止批量执行
- 第一个错误发生时取消上下文
- 剩余任务接收上下文取消
- 适合依赖操作

## 示例

查看 [examples](internal/examples/) 目录:

- [example1](internal/examples/example1) - 基础 errgroup 用法
- [example2](internal/examples/example2) - 批量任务处理
- [example3](internal/examples/example3) - 高级模式

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 许可证类型

MIT 许可证。详见 [LICENSE](LICENSE)。

---

## 🤝 项目贡献

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **发现问题？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **功能建议？** 创建 issue 讨论您的想法
- 📖 **文档疑惑？** 报告问题，帮助我们改进文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，帮助我们优化性能
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **反馈意见？** 欢迎提出建议和意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：为面向用户的更改更新文档，并使用有意义的提交消息
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Merge Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Merge Request 和报告问题来为此项目做出贡献。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**祝你用这个包编程愉快！** 🎉🎉🎉

<!-- TEMPLATE (ZH) END: STANDARD PROJECT FOOTER -->

---

## GitHub 标星点赞

[![Stargazers](https://starchart.cc/yyle88/egobatch.svg?variant=adaptive)](https://starchart.cc/yyle88/egobatch)
