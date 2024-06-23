// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package pgstocksqry

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type OrdersOrderStatus string

const (
	OrdersOrderStatusUnknown         OrdersOrderStatus = "Unknown"
	OrdersOrderStatusNew             OrdersOrderStatus = "New"
	OrdersOrderStatusAwaitingPayment OrdersOrderStatus = "AwaitingPayment"
	OrdersOrderStatusPayed           OrdersOrderStatus = "Payed"
	OrdersOrderStatusCancelled       OrdersOrderStatus = "Cancelled"
	OrdersOrderStatusFailed          OrdersOrderStatus = "Failed"
)

func (e *OrdersOrderStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = OrdersOrderStatus(s)
	case string:
		*e = OrdersOrderStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for OrdersOrderStatus: %T", src)
	}
	return nil
}

type NullOrdersOrderStatus struct {
	OrdersOrderStatus OrdersOrderStatus
	Valid             bool // Valid is true if OrdersOrderStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullOrdersOrderStatus) Scan(value interface{}) error {
	if value == nil {
		ns.OrdersOrderStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.OrdersOrderStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullOrdersOrderStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.OrdersOrderStatus), nil
}

type OrdersOrder struct {
	ID        int64
	UserID    int64
	Status    OrdersOrderStatus
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type OrdersOrderItem struct {
	SkuID     int64
	OrderID   int64
	Count     int64
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type StocksStock struct {
	ID        int64
	Available int64
	Reserved  int64
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}
