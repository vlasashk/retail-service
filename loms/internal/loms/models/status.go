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

func (os OrderStatus) String() string {
	switch os {
	case UnknownStatus:
		return "unknown"
	case NewStatus:
		return "new"
	case AwaitingPaymentStatus:
		return "awaiting payment"
	case PayedStatus:
		return "payed"
	case CancelledStatus:
		return "cancelled"
	case FailedStatus:
		return "failed"
	default:
		return "unknown"
	}
}
