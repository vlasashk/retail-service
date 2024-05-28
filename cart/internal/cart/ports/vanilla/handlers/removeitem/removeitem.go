package removeitem

import (
	"context"
	"net/http"
	"strconv"

	"route256/cart/internal/cart/ports/vanilla/handlers/errhandle"
	"route256/cart/pkg/constants"

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

	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil || userID == 0 {
		localLog.Error().Err(err).Int64("user_id", userID).Send()
		errhandle.NewErr(constants.ErrInvalidUserID).Send(w, localLog, http.StatusBadRequest)
		return
	}

	skuID, err := strconv.ParseInt(r.PathValue("sku_id"), 10, 64)
	if err != nil || skuID == 0 {
		localLog.Error().Err(err).Int64("sku_id", skuID).Send()
		errhandle.NewErr(constants.ErrInvalidSKUID).Send(w, localLog, http.StatusBadRequest)
		return
	}

	if err = h.itemRemover.DeleteItem(r.Context(), userID, skuID); err != nil {
		localLog.Error().Err(err).Send()
		errhandle.NewErr(constants.ErrRemoveItem).Send(w, localLog, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
