package egobatch_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/egobatch"
	"github.com/yyle88/egobatch/erxgroup"
	"github.com/yyle88/egobatch/internal/myassert"
	"github.com/yyle88/egobatch/internal/myerrors"
	"github.com/yyle88/neatjson/neatjsons"
)

func TestTaskOutput(t *testing.T) {
	type Param struct {
		Value int
	}

	type Result struct {
		Value string
	}

	var args []*Param
	for _, num := range []int{0, 1, 2, 3, 4, 5} {
		args = append(args, &Param{Value: num})
	}

	taskBatch := egobatch.NewTaskBatch[*Param, *egobatch.TaskOutput[*Param, *Result, *myerrors.Error], *myerrors.Error](args)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(err error) *myerrors.Error {
		return myerrors.ErrorWrongContext("wrong-ctx. error=%v", err)
	})
	ego := erxgroup.NewGroup[*myerrors.Error](context.Background())
	ego.SetLimit(3)
	taskBatch.EgoRun(ego, func(ctx context.Context, arg *Param) (*egobatch.TaskOutput[*Param, *Result, *myerrors.Error], *myerrors.Error) {
		if arg.Value%3 == 2 {
			return nil, myerrors.ErrorServiceError("wrong-db")
		}
		res := &Result{Value: strconv.Itoa(arg.Value)}
		return egobatch.NewOkTaskOutput[*Param, *Result, *myerrors.Error](arg, res), nil
	})
	myassert.NoError(t, ego.Wait())
	results := taskBatch.Tasks.Flatten(egobatch.NewWaTaskOutput[*Param, *Result, *myerrors.Error])

	ops := egobatch.TaskOutputList[*Param, *Result, *myerrors.Error](results)
	t.Log(neatjsons.S(ops))

	require.Len(t, ops.OkList(), 4)
	require.Len(t, ops.WaList(), 2)

	require.Equal(t, 4, ops.OkCount())
	require.Equal(t, 2, ops.WaCount())

	t.Log(neatjsons.S(ops.OkResults()))
	t.Log(neatjsons.S(ops.WaReasons()))
}
