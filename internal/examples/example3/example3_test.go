package example3_test

import (
	"context"
	"math/rand/v2"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/egobatch"
	"github.com/yyle88/egobatch/erxgroup"
	"github.com/yyle88/egobatch/internal/examples/example3"
	"github.com/yyle88/egobatch/internal/myerrors"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

type Step1Output = egobatch.TaskOutput[*example3.Step1Param, *example3.Step1Result, *myerrors.Error]
type Step2Output = egobatch.TaskOutput[*example3.Step2Param, *example3.Step2Result, *myerrors.Error]
type Step3Output = egobatch.TaskOutput[*example3.Step3Param, *example3.Step3Result, *myerrors.Error]

func TestTaskOutput(t *testing.T) {
	params := example3.NewStep1Params(5)

	outputs := processStep1s(t, params, zaplog.LOGGER)
	t.Log(neatjsons.S(outputs))

	require.Len(t, outputs.OkList(), 3)
	require.Len(t, outputs.WaList(), 2)

	require.Equal(t, 3, outputs.OkCount())
	require.Equal(t, 2, outputs.WaCount())

	t.Log(neatjsons.S(outputs.OkResults()))
	t.Log(neatjsons.S(outputs.WaReasons()))
}

func processStep1s(t *testing.T, params []*example3.Step1Param, zapLog *zaplog.Zap) egobatch.TaskOutputList[*example3.Step1Param, *example3.Step1Result, *myerrors.Error] {
	taskBatch := egobatch.NewTaskBatch[*example3.Step1Param, *Step1Output, *myerrors.Error](params)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(err error) *myerrors.Error {
		return myerrors.ErrorWrongContext("wrong-ctx. error=%v", err)
	})
	ego := erxgroup.NewGroup[*myerrors.Error](context.Background())
	ego.SetLimit(3)
	taskBatch.EgoRun(ego, func(ctx context.Context, arg *example3.Step1Param) (*Step1Output, *myerrors.Error) {
		return processStep1Func(t, ctx, arg, zapLog.SubZap(zap.Int("num_a", arg.NumA)))
	})
	must.Null(ego.Wait())
	outputs := taskBatch.Tasks.Flatten(egobatch.NewWaTaskOutput[*example3.Step1Param, *example3.Step1Result, *myerrors.Error])
	require.Equal(t, len(params), len(outputs))
	return outputs
}

func processStep1Func(t *testing.T, ctx context.Context, arg *example3.Step1Param, zapLog *zaplog.Zap) (*Step1Output, *myerrors.Error) {
	if arg.NumA%2 == 1 {
		zapLog.SUG.Debugln("wrong-a")
		return nil, myerrors.ErrorServiceError("step-1-wrong-db")
	}
	zapLog.SUG.Debugln("process")
	res := &example3.Step1Result{
		ResA:         strconv.Itoa(arg.NumA),
		Step2Outputs: processStep2s(t, example3.NewStep2Params(1+rand.IntN(3)), zapLog),
	}
	return egobatch.NewOkTaskOutput[*example3.Step1Param, *example3.Step1Result, *myerrors.Error](arg, res), nil
}

func processStep2s(t *testing.T, params []*example3.Step2Param, zapLog *zaplog.Zap) egobatch.TaskOutputList[*example3.Step2Param, *example3.Step2Result, *myerrors.Error] {
	taskBatch := egobatch.NewTaskBatch[*example3.Step2Param, *Step2Output, *myerrors.Error](params)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(err error) *myerrors.Error {
		return myerrors.ErrorWrongContext("wrong-ctx. error=%v", err)
	})
	ego := erxgroup.NewGroup[*myerrors.Error](context.Background())
	ego.SetLimit(3)
	taskBatch.EgoRun(ego, func(ctx context.Context, arg *example3.Step2Param) (*Step2Output, *myerrors.Error) {
		return processStep2Func(t, ctx, arg, zapLog.SubZap(zap.Int("num_b", arg.NumB)))
	})
	must.Null(ego.Wait())
	outputs := taskBatch.Tasks.Flatten(egobatch.NewWaTaskOutput[*example3.Step2Param, *example3.Step2Result, *myerrors.Error])
	require.Equal(t, len(params), len(outputs))
	return outputs
}

func processStep2Func(t *testing.T, ctx context.Context, arg *example3.Step2Param, zapLog *zaplog.Zap) (*Step2Output, *myerrors.Error) {
	if rand.IntN(100) < 30 {
		zapLog.SUG.Debugln("wrong-b")
		return nil, myerrors.ErrorServiceError("step-2-wrong-db")
	}
	zapLog.SUG.Debugln("process")
	res := &example3.Step2Result{
		ResB:         strconv.Itoa(arg.NumB),
		Step3Outputs: processStep3s(t, example3.NewStep3Params(1+rand.IntN(3)), zapLog),
	}
	return egobatch.NewOkTaskOutput[*example3.Step2Param, *example3.Step2Result, *myerrors.Error](arg, res), nil
}

func processStep3s(t *testing.T, params []*example3.Step3Param, zapLog *zaplog.Zap) egobatch.TaskOutputList[*example3.Step3Param, *example3.Step3Result, *myerrors.Error] {
	taskBatch := egobatch.NewTaskBatch[*example3.Step3Param, *Step3Output, *myerrors.Error](params)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(err error) *myerrors.Error {
		return myerrors.ErrorWrongContext("wrong-ctx. error=%v", err)
	})
	ego := erxgroup.NewGroup[*myerrors.Error](context.Background())
	ego.SetLimit(3)
	taskBatch.EgoRun(ego, func(ctx context.Context, arg *example3.Step3Param) (*Step3Output, *myerrors.Error) {
		return processStep3Func(t, ctx, arg, zapLog.SubZap(zap.Int("num_c", arg.NumC)))
	})
	must.Null(ego.Wait())
	outputs := taskBatch.Tasks.Flatten(egobatch.NewWaTaskOutput[*example3.Step3Param, *example3.Step3Result, *myerrors.Error])
	require.Equal(t, len(params), len(outputs))
	return outputs
}

func processStep3Func(t *testing.T, ctx context.Context, arg *example3.Step3Param, zapLog *zaplog.Zap) (*Step3Output, *myerrors.Error) {
	if rand.IntN(100) < 50 {
		zapLog.SUG.Debugln("wrong-c")
		return nil, myerrors.ErrorServiceError("step-3-wrong-db")
	}
	zapLog.SUG.Debugln("process")
	res := &example3.Step3Result{
		ResC: strconv.Itoa(arg.NumC),
	}
	return egobatch.NewOkTaskOutput[*example3.Step3Param, *example3.Step3Result, *myerrors.Error](arg, res), nil
}
