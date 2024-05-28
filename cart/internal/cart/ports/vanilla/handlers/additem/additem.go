package additem

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"route256/cart/internal/cart/models"
	"route256/cart/internal/cart/ports/vanilla/handlers/common"
	"route256/cart/internal/cart/ports/vanilla/handlers/errhandle"
	"route256/cart/pkg/constants"

	"github.com/rs/zerolog"
)

//go:generate mockery --name=CartAdder
type CartAdder interface {
	AddItem(ctx context.Context, userID, skuID int64, count uint16) error
}

type Handler struct {
	adder           CartAdder
	productProvider common.ProductProvider
	log             zerolog.Logger
}

func New(log zerolog.Logger, adder CartAdder, provider common.ProductProvider) *Handler {
	return &Handler{
		adder:           adder,
		productProvider: provider,
		log:             log,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	localLog := h.log.With().Str("handler", "add_item").Logger()

	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil || userID <= 0 {
		errhandle.NewErr(constants.ErrInvalidUserID).Send(w, localLog, http.StatusBadRequest)
		return
	}

	skuID, err := strconv.ParseInt(r.PathValue("sku_id"), 10, 64)
	if err != nil || skuID <= 0 {
		errhandle.NewErr(constants.ErrInvalidSKUID).Send(w, localLog, http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(r.Body)
	defer func() {
		if err = r.Body.Close(); err != nil {
			localLog.Error().Err(err).Send()
		}
	}()
	if err != nil {
		errhandle.NewErr(constants.ErrReadBody).Send(w, localLog, http.StatusInternalServerError)
		return
	}
	var count itemCountReq

	err = json.Unmarshal(data, &count)
	if err != nil {
		errhandle.NewErr(constants.ErrUnmarshal).Send(w, localLog, http.StatusInternalServerError)
		return
	}

	if count.Count == 0 {
		errhandle.NewErr(constants.ErrBadCount).Send(w, localLog, http.StatusBadRequest)
		return
	}

	_, err = h.productProvider.GetProduct(r.Context(), skuID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			errhandle.NewErr(constants.ErrItemNotFound).Send(w, localLog, http.StatusPreconditionFailed)
			return
		}
		errhandle.NewErr(constants.ErrGetItem).Send(w, localLog, http.StatusInternalServerError)
		return
	}

	if err = h.adder.AddItem(r.Context(), userID, skuID, count.Count); err != nil {
		errhandle.NewErr(constants.ErrAddItem).Send(w, localLog, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
