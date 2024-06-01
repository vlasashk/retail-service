package removeitem

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

	userID, errUserID := converter.UserToInt(r.PathValue(constants.UserID))
	skuID, errSKUiD := converter.SKUtoInt(r.PathValue(constants.SKUid))
	if err := errors.Join(errUserID, errSKUiD); err != nil {
		localLog.Error().Err(err).Str(constants.UserID, r.PathValue(constants.UserID)).Str(constants.SKUid, r.PathValue(constants.SKUid)).Send()
		errhandle.NewErr(err.Error()).Send(w, localLog, http.StatusBadRequest)
		return
	}

	if err := h.itemRemover.DeleteItem(r.Context(), userID, skuID); err != nil {
		localLog.Error().Err(err).Send()
		errhandle.NewErr(models.ErrRemoveItem.Error()).Send(w, localLog, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
