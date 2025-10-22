package example1_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/yyle88/egobatch"
	"github.com/yyle88/egobatch/erxgroup"
	"github.com/yyle88/egobatch/internal/examples/example1"
	"github.com/yyle88/egobatch/internal/myerrors"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
)

func TestRun(t *testing.T) {
	ctx := context.Background()
	guests := example1.NewGuests(10)
	taskResults := processGuests(ctx, guests)
	//把结果展成平铺的，避免泛型套泛型的输出，这样有利于外部观察和使用
	guestOrdersStates := taskResults.Flatten(func(guest *example1.Guest, err *myerrors.Error) *example1.GuestOrdersStates {
		return &example1.GuestOrdersStates{
			Guest:       guest,
			OrderStates: nil,
			Outline:     "",
			Erx:         err,
		}
	})
	t.Log(neatjsons.S(guestOrdersStates))
}

func processGuests(ctx context.Context, guests []*example1.Guest) egobatch.Tasks[*example1.Guest, *example1.GuestOrdersStates, *myerrors.Error] {
	taskBatch := egobatch.NewTaskBatch[*example1.Guest, *example1.GuestOrdersStates, *myerrors.Error](guests)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(err error) *myerrors.Error {
		return myerrors.ErrorWrongContext("wrong-ctx-can-not-invoke-process-guest-func. error=%v", err)
	})
	ego := erxgroup.NewGroup[*myerrors.Error](ctx)
	ego.SetLimit(3)
	taskBatch.EgoRun(ego, processGuestFunc)
	must.Null(ego.Wait())
	return taskBatch.Tasks
}

func processGuestFunc(ctx context.Context, guest *example1.Guest) (*example1.GuestOrdersStates, *myerrors.Error) {
	if rand.IntN(2) == 0 {
		return nil, myerrors.ErrorServiceError("wrong-db")
	}
	orderCount := 1 + rand.IntN(5)
	orders := example1.NewOrders(guest, orderCount)

	taskResults := processOrders(ctx, orders)

	// 这里把数据降低维度，避免泛型套泛型，能够让逻辑更清楚些，直接返回这个 task-results 也是可以的
	orderStates := taskResults.Flatten(func(order *example1.Order, err *myerrors.Error) *example1.OrderState {
		return &example1.OrderState{
			Order: order,
			Erx:   err,
		}
	})

	outline := createStatusOutline(orderStates)

	return &example1.GuestOrdersStates{
		Guest:       guest,
		OrderStates: orderStates,
		Outline:     outline,
		Erx:         nil,
	}, nil
}

func processOrders(ctx context.Context, orders []*example1.Order) egobatch.Tasks[*example1.Order, *example1.OrderState, *myerrors.Error] {
	taskBatch := egobatch.NewTaskBatch[*example1.Order, *example1.OrderState, *myerrors.Error](orders)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(err error) *myerrors.Error {
		return myerrors.ErrorWrongContext("wrong-ctx-can-not-invoke-process-order-func. error=%v", err)
	})
	ego := erxgroup.NewGroup[*myerrors.Error](ctx)
	ego.SetLimit(2)
	taskBatch.EgoRun(ego, processOrderFunc)
	must.Null(ego.Wait())
	return taskBatch.Tasks
}

func processOrderFunc(ctx context.Context, order *example1.Order) (*example1.OrderState, *myerrors.Error) {
	if rand.IntN(2) == 0 {
		return nil, myerrors.ErrorServiceError("wrong-db")
	}
	return &example1.OrderState{
		Order: order,
		State: "OK",
		Erx:   nil,
	}, nil
}

func createStatusOutline(orderStates []*example1.OrderState) string {
	okCount := 0
	waCount := 0
	for _, state := range orderStates {
		if state.Erx != nil {
			waCount++
		} else {
			okCount++
		}
	}
	return fmt.Sprintf("ok-count:%d wa-count=%d", okCount, waCount)
}
