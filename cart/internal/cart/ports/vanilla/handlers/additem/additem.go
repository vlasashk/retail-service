package additem

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"route256/cart/internal/cart/models"
	"route256/cart/internal/cart/models/constants"
	"route256/cart/internal/cart/ports/vanilla/handlers/common"
	"route256/cart/internal/cart/ports/vanilla/handlers/errhandle"
	"route256/cart/internal/cart/utils/converter"

	"github.com/go-playground/validator/v10"
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

	userID, errUserID := converter.UserToInt(r.PathValue(constants.PathArgUserID))
	skuID, errSKUiD := converter.SKUtoInt(r.PathValue(constants.PathArgSKU))
	if err := errors.Join(errUserID, errSKUiD); err != nil {
		localLog.Error().Err(err).Str(constants.PathArgUserID, r.PathValue(constants.PathArgUserID)).Str(constants.PathArgSKU, r.PathValue(constants.PathArgSKU)).Send()
		errhandle.NewErr(err.Error()).Send(w, localLog, http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(r.Body)
	defer func() {
		if err = r.Body.Close(); err != nil {
			localLog.Error().Err(err).Send()
		}
	}()
	if err != nil {
		localLog.Error().Err(err).Send()
		errhandle.NewErr(constants.ErrReadBody).Send(w, localLog, http.StatusInternalServerError)
		return
	}
	var count itemCountReq

	err = json.Unmarshal(data, &count)
	if err != nil {
		localLog.Error().Err(err).Send()
		errhandle.NewErr(constants.ErrJsonProcessing).Send(w, localLog, http.StatusBadRequest)
		return
	}

	if err = validator.New(validator.WithRequiredStructEnabled()).Struct(count); err != nil {
		localLog.Error().Str("error", constants.ErrBadCount).Send()
		errhandle.NewErr(constants.ErrBadCount).Send(w, localLog, http.StatusBadRequest)
		return
	}

	_, err = h.productProvider.GetProduct(r.Context(), skuID)
	if err != nil {
		localLog.Error().Err(err).Send()
		if errors.Is(err, models.ErrNotFound) {
			errhandle.NewErr(constants.ErrItemNotFound).Send(w, localLog, http.StatusPreconditionFailed)
			return
		}
		errhandle.NewErr(constants.ErrGetItem).Send(w, localLog, http.StatusInternalServerError)
		return
	}

	if err = h.adder.AddItem(r.Context(), userID, skuID, count.Count); err != nil {
		localLog.Error().Err(err).Send()
		errhandle.NewErr(constants.ErrAddItem).Send(w, localLog, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
