package usecase

import (
	"context"

	"route256/loms/internal/loms/models"

	"github.com/rs/zerolog"
)

//go:generate minimock -i stockProvider -p usecase_test
type stockProvider interface {
	Reserve(ctx context.Context, order models.Order) error
	ReserveRemove(ctx context.Context, order models.Order) error
	ReserveCancel(ctx context.Context, order models.Order) error
	GetBySKU(ctx context.Context, skuID uint32) (int64, error)
}

//go:generate minimock -i orderManager -p usecase_test
type orderManager interface {
	Create(ctx context.Context, order models.Order) (int64, error)
	SetStatus(ctx context.Context, orderID int64, status models.OrderStatus) error
	GetByOrderID(ctx context.Context, orderID int64) (models.Order, error)
}

type UseCase struct {
	log    zerolog.Logger
	stocks stockProvider
	orders orderManager
}

func New(
	log zerolog.Logger,
	orders orderManager,
	stocks stockProvider,
) *UseCase {
	return &UseCase{
		log:    log,
		orders: orders,
		stocks: stocks,
	}
}

func (uc *UseCase) OrderCreate(ctx context.Context, order models.Order) (int64, error) {
	var err error

	orderID, err := uc.orders.Create(ctx, order)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			uc.log.Error().Err(err).Str("status", "fail").Int64("orderID", orderID).Msg("reservation failed")
			err = uc.orders.SetStatus(ctx, orderID, models.FailedStatus)
			if err != nil {
				uc.log.Error().Err(err).Str("set fail", "FAILED").Int64("orderID", orderID).Send()
			}
		} else {
			uc.log.Debug().Str("status", "success").Int64("orderID", orderID).Msg("reservation success")
			err = uc.orders.SetStatus(ctx, orderID, models.AwaitingPaymentStatus)
			if err != nil {
				uc.log.Error().Err(err).Str("set fail", "AWAITING_PAYMENT").Int64("orderID", orderID).Send()
			}
		}
	}()

	if err = uc.stocks.Reserve(ctx, order); err != nil {
		return 0, err
	}

	return orderID, nil
}
func (uc *UseCase) OrderInfo(ctx context.Context, orderID int64) (models.Order, error) {
	order, err := uc.orders.GetByOrderID(ctx, orderID)
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}
func (uc *UseCase) OrderPay(ctx context.Context, orderID int64) error {
	order, err := uc.orders.GetByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			uc.log.Error().Err(err).Str("status", "fail").Int64("orderID", orderID).Msg("reservation removal failed")
			err = uc.orders.SetStatus(ctx, orderID, models.FailedStatus)
			if err != nil {
				uc.log.Error().Err(err).Str("set fail", "FAILED").Int64("orderID", orderID).Send()
			}
		} else {
			uc.log.Debug().Str("status", "success").Int64("orderID", orderID).Msg("reservation removal success")
			err = uc.orders.SetStatus(ctx, orderID, models.PayedStatus)
			if err != nil {
				uc.log.Error().Err(err).Str("set fail", "PAYED").Int64("orderID", orderID).Send()
			}
		}
	}()

	return uc.stocks.ReserveRemove(ctx, order)
}
func (uc *UseCase) OrderCancel(ctx context.Context, orderID int64) error {
	order, err := uc.orders.GetByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			uc.log.Error().Err(err).Str("status", "fail").Int64("orderID", orderID).Msg("reservation cancellation failed")
			err = uc.orders.SetStatus(ctx, orderID, models.FailedStatus)
			if err != nil {
				uc.log.Error().Err(err).Str("set fail", "FAILED").Int64("orderID", orderID).Send()
			}
		} else {
			uc.log.Debug().Str("status", "success").Int64("orderID", orderID).Msg("reservation cancellation success")
			err = uc.orders.SetStatus(ctx, orderID, models.CancelledStatus)
			if err != nil {
				uc.log.Error().Err(err).Str("set fail", "CANCELLED").Int64("orderID", orderID).Send()
			}
		}
	}()

	return uc.stocks.ReserveCancel(ctx, order)
}
func (uc *UseCase) StocksInfo(ctx context.Context, skuID uint32) (int64, error) {
	amount, err := uc.stocks.GetBySKU(ctx, skuID)
	if err != nil {
		return 0, err
	}

	return amount, nil
}
