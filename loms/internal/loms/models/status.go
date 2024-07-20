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

func (s OrderStatus) String() string {
	switch s {
	case NewStatus:
		return "New"
	case AwaitingPaymentStatus:
		return "AwaitingPayment"
	case PayedStatus:
		return "Payed"
	case CancelledStatus:
		return "Cancelled"
	case FailedStatus:
		return "Failed"
	default:
		return "Unknown"
	}
}
