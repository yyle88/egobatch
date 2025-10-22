[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yyle88/egobatch/release.yml?branch=main&label=BUILD)](https://github.com/yyle88/egobatch/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yyle88/egobatch)](https://pkg.go.dev/github.com/yyle88/egobatch)
[![Coverage Status](https://img.shields.io/coveralls/github/yyle88/egobatch/main.svg)](https://coveralls.io/github/yyle88/egobatch?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25%2B-lightgrey.svg)](https://github.com/yyle88/egobatch)
[![GitHub Release](https://img.shields.io/github/release/yyle88/egobatch.svg)](https://github.com/yyle88/egobatch/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yyle88/egobatch)](https://goreportcard.com/report/github.com/yyle88/egobatch)

# egobatch

Generic batch task processing with type-safe custom error handling.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->
## CHINESE README

[中文说明](README.zh.md)
<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## Main Features

🎯 **Type-Safe Error Handling**: Generic error types with comparable constraints
⚡ **Batch Task Processing**: Concurrent execution with TaskBatch abstraction
🔄 **Flexible Execution Modes**: Glide mode (independent tasks) and fail-fast mode
🌍 **Context Cancellation**: Full context propagation and timeout support
📋 **Result Aggregation**: Filter success/failed tasks with OkTasks/WaTasks methods

## Installation

```bash
go get github.com/yyle88/egobatch
```

## Quick Start

### Basic errgroup with Custom Error Type

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/yyle88/egobatch/erxgroup"
)

// MyError is a custom error type with Code and Msg fields
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
	ego.Go(func(ctx context.Context) *MyError {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("Task 1 finished OK")
		return nil
	})

	// Add task 2: takes 50ms to finish
	ego.Go(func(ctx context.Context) *MyError {
		time.Sleep(50 * time.Millisecond)
		fmt.Println("Task 2 finished OK")
		return nil
	})

	// Add task 3: takes 80ms to finish
	ego.Go(func(ctx context.Context) *MyError {
		time.Sleep(80 * time.Millisecond)
		fmt.Println("Task 3 finished OK")
		return nil
	})

	// Wait until tasks finish and get the first issue
	if erx := ego.Wait(); erx != nil {
		fmt.Printf("Got issue: %s\n", erx.Error())
	} else {
		fmt.Println("Tasks finished OK")
	}
}
```

⬆️ **Source:** [Source](internal/demos/demo1x/main.go)

### Batch Task Processing

```go
package main

import (
	"context"
	"fmt"

	"github.com/yyle88/egobatch"
	"github.com/yyle88/egobatch/erxgroup"
	"github.com/yyle88/must"
)

// MyError is a simple custom error type with error code
type MyError struct {
	Code string
}

func (e *MyError) Error() string {
	return e.Code
}

func main() {
	// Create batch with arguments
	args := []int{1, 2, 3, 4, 5}
	batch := egobatch.NewTaskBatch[int, string, *MyError](args)

	// Configure glide mode - keep going even when issues happen
	batch.SetGlide(true)

	// Execute batch tasks
	ctx := context.Background()
	ego := erxgroup.NewGroup[*MyError](ctx)

	batch.EgoRun(ego, func(ctx context.Context, num int) (string, *MyError) {
		if num%2 == 0 {
			// Even numbers finish OK
			return fmt.Sprintf("even-%d", num), nil
		}
		// Odd numbers have issues
		return "", &MyError{Code: "ODD_NUMBER"}
	})

	// In glide mode, ego.Wait() returns nil because errors are captured in tasks
	must.Null(ego.Wait())

	// Get and handle task outcomes
	okTasks := batch.Tasks.OkTasks()
	waTasks := batch.Tasks.WaTasks()

	fmt.Printf("Success: %d, Failed: %d\n", len(okTasks), len(waTasks))

	// Show OK outcomes
	for _, task := range okTasks {
		fmt.Printf("Arg: %d -> Outcome: %s\n", task.Arg, task.Res)
	}

	// Show bad outcomes
	for _, task := range waTasks {
		fmt.Printf("Arg: %d -> Issue: %s\n", task.Arg, task.Erx.Error())
	}
}
```

⬆️ **Source:** [Source](internal/demos/demo2x/main.go)

### Context Timeout Handling

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

// MyError is a custom error type with error code
type MyError struct {
	Code string
}

func (e *MyError) Error() string {
	return e.Code
}

func main() {
	// Create context with 150ms timeout
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	// Create batch with task arguments
	args := []int{1, 2, 3, 4, 5}
	batch := egobatch.NewTaskBatch[int, string, *MyError](args)

	// Use glide mode to see which tasks finish vs timeout
	batch.SetGlide(true)

	// Convert context issues to custom error type
	batch.SetWaCtx(func(err error) *MyError {
		return &MyError{Code: "TIMEOUT"}
	})

	ego := erxgroup.NewGroup[*MyError](ctx)

	batch.EgoRun(ego, func(ctx context.Context, num int) (string, *MyError) {
		// Each task needs different time: 50ms, 100ms, 150ms, 200ms, 250ms
		taskTime := time.Duration(num*50) * time.Millisecond

		timer := time.NewTimer(taskTime)
		defer timer.Stop()

		select {
		case <-timer.C:
			// Task finishes within timeout
			fmt.Printf("Task %d finished (%dms)\n", num, num*50)
			return fmt.Sprintf("task-%d", num), nil
		case <-ctx.Done():
			// Task cancelled due to timeout
			fmt.Printf("Task %d cancelled (%dms needed)\n", num, num*50)
			return "", &MyError{Code: "CANCELLED"}
		}
	})

	// In glide mode, ego.Wait() returns nil because errors are captured in tasks
	must.Null(ego.Wait())

	// Show task outcomes
	okTasks := batch.Tasks.OkTasks()
	waTasks := batch.Tasks.WaTasks()

	fmt.Printf("\nSuccess: %d, Timeout: %d\n", len(okTasks), len(waTasks))

	// Show finished tasks
	for _, task := range okTasks {
		fmt.Printf("Arg: %d -> Outcome: %s\n", task.Arg, task.Res)
	}

	// Show timed-out tasks
	for _, task := range waTasks {
		fmt.Printf("Arg: %d -> Issue: %s\n", task.Arg, task.Erx.Error())
	}
}
```

⬆️ **Source:** [Source](internal/demos/demo3x/main.go)

### Fail-Fast Mode

```go
batch := egobatch.NewTaskBatch[int, string, *MyError](args)
// Default is fail-fast mode (Glide: false)

ego := erxgroup.NewGroup[*MyError](ctx)
batch.EgoRun(ego, taskFunc)

if erx := ego.Wait(); erx != nil {
    // First error stops execution
    fmt.Printf("Stopped on error: %s\n", erx.Error())
}
```

### Task Result Transformation

```go
tasks := batch.Tasks

// Flatten with error handling
results := tasks.Flatten(func(arg int, err *MyError) string {
    return fmt.Sprintf("error-%d: %s", arg, err.Code)
})

// Mix of success results and transformed errors
for _, result := range results {
    fmt.Println(result)
}
```

## Core Components

### erxgroup.Group[E ErrorType]

Generic wrapper around `errgroup.Group` with type-safe custom errors:

- `NewGroup[E](ctx)`: Create new group with custom error type
- `Go(func(ctx) E)`: Add task returning custom error
- `TryGo(func(ctx) E)`: Add task with limit checking
- `Wait() E`: Wait and get first typed error
- `SetLimit(n)`: Restrict concurrent task count

### TaskBatch[A, R, E]

Batch task execution with concurrent processing:

- `NewTaskBatch[A, R, E](args)`: Create batch from arguments
- `SetGlide(bool)`: Configure execution mode
- `SetWaCtx(func(error) E)`: Handle context errors
- `GetRun(idx, func)`: Get task execution function
- `EgoRun(ego, func)`: Run batch with errgroup

### Tasks[A, R, E]

Task collection with filtering methods:

- `OkTasks()`: Get tasks that completed with success
- `WaTasks()`: Get tasks that failed
- `Flatten(func)`: Transform results with error handling

## Advanced Usage

### Context Timeout Handling

```go
batch := egobatch.NewTaskBatch[int, string, *MyError](args)
batch.SetGlide(true)

// Convert context errors to custom type
batch.SetWaCtx(func(err error) *MyError {
    return &MyError{Code: "CONTEXT_ERROR", Msg: err.Error()}
})

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

ego := erxgroup.NewGroup[*MyError](ctx)
batch.EgoRun(ego, taskFunc)
ego.Wait()

// Tasks that timed out have context errors recorded
for _, task := range batch.Tasks.WaTasks() {
    fmt.Printf("Task %d error: %s\n", task.Arg, task.Erx.Error())
}
```

### TaskOutput Pattern

```go
import "github.com/yyle88/egobatch"

// Create task outputs
outputs := egobatch.TaskOutputList[int, string, *MyError]{
    egobatch.NewOkTaskOutput(1, "success-1"),
    egobatch.NewWaTaskOutput(2, &MyError{Code: "FAIL"}),
    egobatch.NewOkTaskOutput(3, "success-3"),
}

// Filter and aggregate
okList := outputs.OkList()
okCount := outputs.OkCount()
results := outputs.OkResults()
errors := outputs.WaReasons()

fmt.Printf("Success count: %d\n", okCount)
fmt.Printf("Results: %v\n", results)
```

## Design Patterns

### ErrorType Constraint

Custom error types must satisfy the `ErrorType` constraint:

```go
type ErrorType interface {
    error
    comparable
}
```

This enables:
- Type-safe error checking with `errors.Is`
- Zero-value nil detection with `constraint.Pass(erx)`
- Custom error conversion with `errors.As`

### Glide vs Fail-Fast

**Glide Mode (Glide: true)**:
- Tasks execute in independent mode
- Errors recorded without stopping others
- Context cancellation affects remaining tasks
- Best with independent operations

**Fail-Fast Mode (Glide: false)**:
- First error stops batch execution
- Context gets cancelled on first error
- Remaining tasks receive context cancellation
- Best with dependent operations

## Examples

See the [examples](internal/examples/) directory:

- [example1](internal/examples/example1) - Basic errgroup usage
- [example2](internal/examples/example2) - Batch task processing
- [example3](internal/examples/example3) - Advanced patterns

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 License

MIT License. See [LICENSE](LICENSE).

---

## 🤝 Contributing

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Found a mistake?** Open an issue on GitHub with reproduction steps
- 💡 **Have a feature idea?** Create an issue to discuss the suggestion
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share the use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize through reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo to get new releases and features
- 🌟 **Success stories?** Share how this package improved the workflow
- 💬 **Feedback?** We welcome suggestions and comments

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage UI).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement the changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation to support client-facing changes and use significant commit messages
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a merge request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project via submitting merge requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Have Fun Coding with this package!** 🎉🎉🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/yyle88/egobatch.svg?variant=adaptive)](https://starchart.cc/yyle88/egobatch)
