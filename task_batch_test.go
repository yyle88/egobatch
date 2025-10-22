package egobatch_test

import (
	"context"
	"math/rand/v2"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/egobatch"
	"github.com/yyle88/egobatch/erxgroup"
	"github.com/yyle88/egobatch/internal/myassert"
	"github.com/yyle88/egobatch/internal/myerrors"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

func TestGroup_Go_TaskRun(t *testing.T) {
	ego := erxgroup.NewGroup[*myerrors.Error](context.Background())
	ego.SetLimit(10)

	args := make([]uint64, 0, 50)
	for num := uint64(0); num < 50; num++ {
		args = append(args, num)
	}

	var taskBatch = egobatch.NewTaskBatch[uint64, string, *myerrors.Error](args)
	for idx := 0; idx < 50; idx++ {
		ego.Go(taskBatch.GetRun(idx, taskRun))
	}
	myassert.Error(t, ego.Wait())

	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erx:", task.Erx)
	}
}

func taskRun(ctx context.Context, arg uint64) (string, *myerrors.Error) {
	if ctx.Err() != nil {
		zaplog.LOG.Info("task no", zap.Uint64("arg", arg))
		return "", myerrors.ErrorWrongContext("error=%v", ctx.Err())
	}
	time.Sleep(time.Duration(rand.IntN(1000)) * time.Millisecond) // 模拟计算延迟
	if arg%10 == 3 {
		zaplog.LOG.Info("task wa", zap.Uint64("arg", arg))
		return "", myerrors.ErrorServiceError("task wa %d", arg) // 模拟某个任务失败
	}
	zaplog.LOG.Info("task ok", zap.Uint64("arg", arg))

	res := strconv.FormatUint(arg, 10)
	return res, nil
}

func TestGroup_Go_SetGlide_TaskRun(t *testing.T) {
	ego := erxgroup.NewGroup[*myerrors.Error](context.Background())
	ego.SetLimit(10)

	args := make([]uint64, 0, 50)
	for num := uint64(0); num < 50; num++ {
		args = append(args, num)
	}

	taskBatch := egobatch.NewTaskBatch[uint64, string, *myerrors.Error](args)
	taskBatch.SetGlide(true)
	for idx := 0; idx < 50; idx++ {
		ego.Go(taskBatch.GetRun(idx, taskRun))
	}
	myassert.NoError(t, ego.Wait())

	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erx:", task.Erx)
	}
}

func TestGroup_Go_SetGlide_SetWaCtx_TaskRun(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*20)
	defer cancelFunc()

	ego := erxgroup.NewGroup[*myerrors.Error](ctx)
	ego.SetLimit(10)

	args := make([]uint64, 0, 50)
	for num := uint64(0); num < 50; num++ {
		args = append(args, num)
	}

	taskBatch := egobatch.NewTaskBatch[uint64, string, *myerrors.Error](args)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(err error) *myerrors.Error {
		return myerrors.ErrorWrongContext("ctx wrong reason=%v", err)
	})
	for idx := 0; idx < 50; idx++ {
		ego.Go(taskBatch.GetRun(idx, func(ctx context.Context, arg uint64) (string, *myerrors.Error) {
			time.Sleep(time.Millisecond * 10)
			res := strconv.FormatUint(arg, 10)
			return res, nil
		}))
	}
	myassert.NoError(t, ego.Wait())

	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erx:", task.Erx)
	}
}

func TestTaskBatch_GetRun(t *testing.T) {
	var args = []uint64{0, 1, 2, 3, 4, 5}
	var taskBatch = egobatch.NewTaskBatch[uint64, string, *myerrors.Error](args)
	for idx, task := range taskBatch.Tasks {
		require.Equal(t, idx, int(task.Arg))
	}

	ctx := context.Background()
	for idx := 0; idx < len(args); idx++ {
		run := taskBatch.GetRun(idx, func(ctx context.Context, arg uint64) (string, *myerrors.Error) {
			res := strconv.FormatUint(arg, 10)
			return res, nil
		})
		erk := run(ctx)
		t.Log(erk)
		myassert.NoError(t, erk)
	}
	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erx:", task.Erx)
		require.Equal(t, strconv.Itoa(idx), task.Res)
		myassert.NoError(t, task.Erx)
	}
	results := taskBatch.Tasks.Flatten(func(arg uint64, erk *myerrors.Error) string {
		return "wa-" + strconv.Itoa(int(arg))
	})
	t.Log(neatjsons.S(results))
	require.Equal(t, []string{"0", "1", "2", "3", "4", "5"}, results)
}

func TestTaskBatch_SetGlide_GetRun(t *testing.T) {
	var args = []uint64{0, 1, 2, 3, 4, 5}
	taskBatch := egobatch.NewTaskBatch[uint64, string, *myerrors.Error](args)
	for idx, task := range taskBatch.Tasks {
		require.Equal(t, idx, int(task.Arg))
	}
	taskBatch.SetGlide(true)

	ctx := context.Background()
	for idx := 0; idx < len(args); idx++ {
		run := taskBatch.GetRun(idx, func(ctx context.Context, arg uint64) (string, *myerrors.Error) {
			if arg%2 == 0 {
				return "", myerrors.ErrorServiceError("wrong db")
			}
			res := strconv.FormatUint(arg, 10)
			return res, nil
		})
		erk := run(ctx)
		t.Log(erk)
		myassert.NoError(t, erk) //当设置 "平滑继续" 时这里不返回错误
	}
	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erx:", task.Erx)
		if idx%2 == 0 {
			require.True(t, myerrors.IsServiceError(task.Erx))
		} else {
			require.Equal(t, strconv.Itoa(idx), task.Res)
			myassert.NoError(t, task.Erx)
		}
	}
	results := taskBatch.Tasks.Flatten(func(arg uint64, erk *myerrors.Error) string {
		return "wa-" + strconv.Itoa(int(arg))
	})
	t.Log(neatjsons.S(results))
	require.Equal(t, []string{"wa-0", "1", "wa-2", "3", "wa-4", "5"}, results)
}

func TestTaskBatch_EgoRun(t *testing.T) {
	ctx := context.Background()

	var args = []uint64{0, 1, 2, 3, 4, 5}
	taskBatch := egobatch.NewTaskBatch[uint64, string, *myerrors.Error](args)
	for idx, task := range taskBatch.Tasks {
		require.Equal(t, idx, int(task.Arg))
	}
	taskBatch.SetGlide(true)

	ego := erxgroup.NewGroup[*myerrors.Error](ctx)
	ego.SetLimit(3)
	taskBatch.EgoRun(ego, func(ctx context.Context, arg uint64) (string, *myerrors.Error) {
		if arg%2 == 0 {
			return "", myerrors.ErrorServiceError("wrong db")
		}
		res := strconv.FormatUint(arg, 10)
		return res, nil
	})
	myassert.NoError(t, ego.Wait())

	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erx:", task.Erx)
		if idx%2 == 0 {
			require.True(t, myerrors.IsServiceError(task.Erx))
		} else {
			require.Equal(t, strconv.Itoa(idx), task.Res)
			myassert.NoError(t, task.Erx)
		}
	}
	results := taskBatch.Tasks.Flatten(func(arg uint64, erk *myerrors.Error) string {
		return "wa-" + strconv.Itoa(int(arg))
	})
	t.Log(neatjsons.S(results))
	require.Equal(t, []string{"wa-0", "1", "wa-2", "3", "wa-4", "5"}, results)
}
