// Package example1 demonstrates basic batch task processing patterns
// Shows guest and order processing with error handling
//
// 包 example1 演示基础批量任务处理模式
// 展示访客和订单处理以及错误处理
package example1

import (
	"fmt"

	"github.com/yyle88/egobatch/internal/myerrors"
)

// Guest represents a guest entity with name and ID
// 访客实体，包含名称和ID
type Guest struct {
	Name    string // Guest name // 访客名称
	GuestID int    `json:"-"` // Guest unique ID // 访客唯一ID
}

// GuestOrdersStates aggregates guest with associated order states
// 聚合访客及其关联的订单状态
type GuestOrdersStates struct {
	Guest       *Guest          // Guest reference // 访客引用
	OrderStates []*OrderState   // Collection of order states // 订单状态集合
	Outline     string          // Summary outline text // 概要文本
	Erx         *myerrors.Error // Processing error if any // 处理错误（如有）
}

// Order represents an order entity with name and ID
// 订单实体，包含名称和ID
type Order struct {
	Name    string // Order name // 订单名称
	OrderID int    `json:"-"` // Order unique ID // 订单唯一ID
}

// OrderState represents order processing state with error tracking
// 订单处理状态，包含错误跟踪
type OrderState struct {
	Order *Order          // Order reference // 订单引用
	State string          // Current processing state // 当前处理状态
	Erx   *myerrors.Error // Processing error if any // 处理错误（如有）
}

// NewGuests creates a collection of guest instances
// Assigns sequential ID and formatted name
//
// NewGuests 创建访客实例集合
// ID和名称按索引生成
func NewGuests(guestCount int) []*Guest {
	var guests = make([]*Guest, 0, guestCount)
	for idx := 0; idx < guestCount; idx++ {
		guests = append(guests, &Guest{
			Name:    fmt.Sprintf("guest(%d)", idx),
			GuestID: idx,
		})
	}
	return guests
}

// NewOrders creates a collection of orders belonging to a guest
// Assigns sequential ID and formatted name with guest reference
//
// NewOrders 创建属于某个访客的订单集合
// ID和名称按索引生成（名称含访客信息）
func NewOrders(guest *Guest, orderCount int) []*Order {
	orders := make([]*Order, 0, orderCount)
	for idx := 0; idx < orderCount; idx++ {
		orders = append(orders, &Order{
			Name:    fmt.Sprintf("guest(%d) order(%d)", guest.GuestID, idx),
			OrderID: idx,
		})
	}
	return orders
}
