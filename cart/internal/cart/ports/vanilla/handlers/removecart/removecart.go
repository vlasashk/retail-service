package removecart

import (
	"context"
	"errors"
	"net/http"

	"route256/cart/internal/cart/models"
	"route256/cart/internal/cart/models/constants"
	"route256/cart/internal/cart/ports/vanilla/handlers/errhandle"
	"route256/cart/internal/cart/utils/converter"

	"github.com/rs/zerolog"
)

//go:generate mockery --name=CartRemover
type CartRemover interface {
	DeleteItemsByUserID(ctx context.Context, userID int64) error
}
type Handler struct {
	cartRemover CartRemover
	log         zerolog.Logger
}

func New(log zerolog.Logger, remover CartRemover) *Handler {
	return &Handler{
		cartRemover: remover,
		log:         log,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	localLog := h.log.With().Str("handler", "remove_cart").Logger()

	userID, err := converter.UserToInt(r.PathValue(constants.PathArgUserID))
	if err != nil {
		localLog.Error().Err(err).Str(constants.PathArgUserID, r.PathValue(constants.PathArgUserID)).Send()
		errhandle.NewErr(err.Error()).Send(w, localLog, http.StatusBadRequest)
		return
	}

	if err = h.cartRemover.DeleteItemsByUserID(r.Context(), userID); err != nil && !errors.Is(err, models.ErrCartIsEmpty) {
		localLog.Error().Err(err).Send()
		errhandle.NewErr(constants.ErrRemoveCart).Send(w, localLog, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
