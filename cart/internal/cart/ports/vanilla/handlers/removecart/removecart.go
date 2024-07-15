package removecart

import (
	"context"
	"errors"
	"net/http"

	"route256/cart/internal/cart/constants"
	"route256/cart/internal/cart/models"
	"route256/cart/internal/cart/ports/vanilla/handlers/errhandle"
	"route256/cart/internal/cart/utils/converter"

	"github.com/rs/zerolog"
)

//go:generate minimock -i CartRemover -p removecart_test
type CartRemover interface {
	DeleteItemsByUserID(ctx context.Context, userID int64) error
}
type Handler struct {
	cartRemover CartRemover
}

func New(remover CartRemover) *Handler {
	return &Handler{
		cartRemover: remover,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	localLog := zerolog.Ctx(r.Context()).With().Str("handler", "remove_cart").Logger()

	userID, err := converter.UserToInt(r.PathValue(constants.UserID))
	if err != nil {
		localLog.Error().Err(err).Str(constants.UserID, r.PathValue(constants.UserID)).Send()
		errhandle.NewErr(err.Error()).Send(w, localLog, http.StatusBadRequest)
		return
	}

	if err = h.cartRemover.DeleteItemsByUserID(r.Context(), userID); err != nil && !errors.Is(err, models.ErrCartIsEmpty) {
		localLog.Error().Err(err).Send()
		errhandle.NewErr(models.ErrRemoveCart.Error()).Send(w, localLog, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
