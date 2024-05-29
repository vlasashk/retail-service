package removeitem

import (
	"context"
	"errors"
	"net/http"

	"route256/cart/internal/cart/models/constants"
	"route256/cart/internal/cart/ports/vanilla/handlers/errhandle"
	"route256/cart/internal/cart/utils/converter"

	"github.com/rs/zerolog"
)

//go:generate mockery --name=ItemRemover
type ItemRemover interface {
	DeleteItem(ctx context.Context, userID, skuID int64) error
}

type Handler struct {
	itemRemover ItemRemover
	log         zerolog.Logger
}

func New(log zerolog.Logger, remover ItemRemover) *Handler {
	return &Handler{
		itemRemover: remover,
		log:         log,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	localLog := h.log.With().Str("handler", "remove_item").Logger()

	userID, errUserID := converter.UserToInt(r.PathValue(constants.PathArgUserID))
	skuID, errSKUiD := converter.SKUtoInt(r.PathValue(constants.PathArgSKU))
	if err := errors.Join(errUserID, errSKUiD); err != nil {
		localLog.Error().Err(err).Str(constants.PathArgUserID, r.PathValue(constants.PathArgUserID)).Str(constants.PathArgSKU, r.PathValue(constants.PathArgSKU)).Send()
		errhandle.NewErr(err.Error()).Send(w, localLog, http.StatusBadRequest)
		return
	}

	if err := h.itemRemover.DeleteItem(r.Context(), userID, skuID); err != nil {
		localLog.Error().Err(err).Send()
		errhandle.NewErr(constants.ErrRemoveItem).Send(w, localLog, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
