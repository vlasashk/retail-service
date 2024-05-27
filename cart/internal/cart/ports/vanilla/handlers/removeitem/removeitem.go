package removeitem

import (
	"context"
	"net/http"
	"strconv"

	"route256/cart/internal/cart/ports/vanilla/handlers/errhandle"

	"github.com/rs/zerolog"
)

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
		errhandle.NewErr("invalid user_id").Send(w, localLog, http.StatusBadRequest)
		return
	}

	skuID, err := strconv.ParseInt(r.PathValue("sku_id"), 10, 64)
	if err != nil || skuID == 0 {
		errhandle.NewErr("invalid sku_id").Send(w, localLog, http.StatusBadRequest)
		return
	}

	if err = h.itemRemover.DeleteItem(r.Context(), userID, skuID); err != nil {
		errhandle.NewErr("failed to remove item").Send(w, localLog, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
