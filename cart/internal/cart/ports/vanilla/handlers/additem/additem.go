package additem

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"route256/cart/internal/cart/constants"
	"route256/cart/internal/cart/models"
	"route256/cart/internal/cart/ports/vanilla/handlers/errhandle"
	"route256/cart/internal/cart/utils/converter"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

//go:generate minimock -i CartAdder -p additem_test
type CartAdder interface {
	AddItem(ctx context.Context, userID, skuID int64, count uint16) error
}

type Handler struct {
	adder CartAdder
}

func New(adder CartAdder) *Handler {
	return &Handler{
		adder: adder,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	localLog := zerolog.Ctx(r.Context()).With().Str("handler", "add_item").Logger()

	// Логирование внутри парсера
	reqData, err := parseDataFromReq(localLog, r)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidSKUID) || errors.Is(err, models.ErrInvalidUserID) ||
			errors.Is(err, models.ErrJSONProcessing) || errors.Is(err, models.ErrBadCount):
			errhandle.NewErr(err.Error()).Send(w, localLog, http.StatusBadRequest)
		case errors.Is(err, models.ErrReadBody):
			errhandle.NewErr(err.Error()).Send(w, localLog, http.StatusInternalServerError)
		default:
			errhandle.NewErr(models.ErrAddItem.Error()).Send(w, localLog, http.StatusInternalServerError)
		}
		return
	}

	if err = h.adder.AddItem(r.Context(), reqData.UserID, reqData.SKUid, reqData.Count); err != nil {
		localLog.Error().Err(err).Send()
		switch {
		case errors.Is(err, models.ErrNotFound):
			errhandle.NewErr(models.ErrNotFound.Error()).Send(w, localLog, http.StatusPreconditionFailed)
		case errors.Is(err, models.ErrInsufficientStock):
			errhandle.NewErr(models.ErrInsufficientStock.Error()).Send(w, localLog, http.StatusPreconditionFailed)
		case errors.Is(err, models.ErrItemProvider):
			errhandle.NewErr(models.ErrItemProvider.Error()).Send(w, localLog, http.StatusInternalServerError)
		case errors.Is(err, models.ErrStockProvider):
			errhandle.NewErr(models.ErrStockProvider.Error()).Send(w, localLog, http.StatusInternalServerError)
		default:
			errhandle.NewErr(models.ErrAddItem.Error()).Send(w, localLog, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func parseDataFromReq(log zerolog.Logger, r *http.Request) (itemCountReq, error) {
	var itemData itemCountReq

	data, err := io.ReadAll(r.Body)
	defer func() {
		if err = r.Body.Close(); err != nil {
			log.Error().Err(err).Send()
		}
	}()
	if err != nil {
		log.Error().Err(err).Send()
		return itemCountReq{}, models.ErrReadBody
	}

	err = json.Unmarshal(data, &itemData)
	if err != nil {
		log.Error().Err(err).Send()
		return itemCountReq{}, models.ErrJSONProcessing
	}

	userID, errUserID := converter.UserToInt(r.PathValue(constants.UserID))
	skuID, errSKUiD := converter.SKUtoInt(r.PathValue(constants.SKUid))
	if err = errors.Join(errUserID, errSKUiD); err != nil {
		log.Error().Err(err).Str(constants.UserID, r.PathValue(constants.UserID)).Send()
		return itemCountReq{}, err
	}

	if err = validator.New(validator.WithRequiredStructEnabled()).Struct(itemData); err != nil {
		log.Error().Err(err).Send()
		return itemCountReq{}, models.ErrBadCount
	}

	itemData.UserID = userID
	itemData.SKUid = skuID

	return itemData, nil
}
