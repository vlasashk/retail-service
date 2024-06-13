package impl

import (
	"errors"

	"route256/loms/internal/loms/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	switch {
	case errors.Is(err, models.ErrOrderNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, models.ErrItemNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, models.ErrInsufficientStock):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, models.ErrReservationConflict):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, models.ErrPaymentStatusConflict):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
