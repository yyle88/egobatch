package erxgroup_test

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/egobatch/erxgroup"
	"github.com/yyle88/egobatch/internal/myassert"
	"github.com/yyle88/egobatch/internal/myerrors"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func TestGoErrGroup(t *testing.T) {
	ctx := context.Background()

	ego, ctx := errgroup.WithContext(ctx) //使用同名的 ctx 覆盖旧的 ctx 这是 go 里面的习惯

	for idx := 0; idx < 50; idx++ {
		num := idx
		ego.Go(func() error {
			t.Log(num)
			return nil
		})
	}

	// 在 wait 里面会调用 cancelFunc 把 ctx 取消掉
	require.NoError(t, ego.Wait())
	// 这里 ctx 报已取消，因为前面 ctx 被覆盖，而覆盖ctx又是go语言开发者的习惯（就是说 errgroup 设计的不太符合习惯，需要改改）
	t.Log("ctx-err-res:", ctx.Err())
	// 这里其实是不符合预期的，因为 ctx 还要被后续逻辑用到
	// 这样就容易导致BUG，因此我在该项目里使用 NewGroup 封装 errgroup.WithContext 把 ctx 隐藏起来，具体请看下面的测试用例
	require.ErrorIs(t, checkCtx(ctx), context.Canceled)
}

func TestNewGroup(t *testing.T) {
	ctx := context.Background()

	ego := erxgroup.NewGroup[*myerrors.Error](ctx)

	for idx := 0; idx < 50; idx++ {
		num := idx
		ego.Go(func(ctx context.Context) *myerrors.Error {
			t.Log(num)
			return nil
		})
	}

	// 在 wait 里面会调用 cancelFunc 把 ctx 取消掉，但是取消的是内部的 ctx 而不是外部的
	myassert.NoError(t, ego.Wait())
	// 这里不受影响
	t.Log("ctx-err-res:", ctx.Err())
	// 这里依然可以用 ctx， 因为它是最外层的 ctx，其不受内部的 cancelFunc 的影响
	// 这样不容易出BUG，但由于 group 的 ctx 被隐藏，group 的 Go 和 TryGo 的 run 都需要是带有 ctx 信息参数的
	require.NoError(t, checkCtx(ctx))
}

func checkCtx(ctx context.Context) error {
	return ctx.Err()
}

func TestNewGroup_StepRun(t *testing.T) {
	ego := erxgroup.NewGroup[*myerrors.Error](context.Background())
	ego.SetLimit(10)

	for idx := 0; idx < 50; idx++ {
		num := idx
		ego.Go(func(ctx context.Context) *myerrors.Error {
			return stepRun(ctx, num)
		})
	}

	myassert.Error(t, ego.Wait())
}

func stepRun(ctx context.Context, idx int) *myerrors.Error {
	if ctx.Err() != nil {
		zaplog.LOG.Info("task no", zap.Int("num", idx))
		return myerrors.ErrorWrongContext("error=%v", ctx.Err())
	}
	time.Sleep(time.Duration(rand.IntN(1000)) * time.Millisecond) // 模拟计算延迟
	if idx%10 == 3 {
		zaplog.LOG.Info("task wa", zap.Int("num", idx))
		return myerrors.ErrorServiceError("task wa %d", idx) // 模拟某个任务失败
	}
	zaplog.LOG.Info("task ok", zap.Int("num", idx))
	return nil
}
