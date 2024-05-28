package removecart

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"route256/cart/internal/cart/models"
	"route256/cart/internal/cart/ports/vanilla/handlers/errhandle"
	"route256/cart/pkg/constants"

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

	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil || userID == 0 {
		errhandle.NewErr(constants.ErrInvalidUserID).Send(w, localLog, http.StatusBadRequest)
		return
	}

	if err = h.cartRemover.DeleteItemsByUserID(r.Context(), userID); err != nil && !errors.Is(err, models.ErrCartIsEmpty) {
		errhandle.NewErr(constants.ErrRemoveCart).Send(w, localLog, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
