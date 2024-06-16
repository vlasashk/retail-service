package models

type OrderStatus uint8

const (
	UnknownStatus OrderStatus = iota
	NewStatus
	AwaitingPaymentStatus
	PayedStatus
	CancelledStatus
	FailedStatus
)
