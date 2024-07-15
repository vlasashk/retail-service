package getcart

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"route256/cart/internal/cart/constants"
	"route256/cart/internal/cart/models"
	"route256/cart/internal/cart/ports/vanilla/handlers/errhandle"
	"route256/cart/internal/cart/utils/converter"

	"github.com/rs/zerolog"
)

//go:generate minimock -i CartRetriever -p getcart_test
type CartRetriever interface {
	GetItemsByUserID(ctx context.Context, userID int64) (models.ItemsInCart, error)
}

type Handler struct {
	retriever CartRetriever
}

func New(retriever CartRetriever) *Handler {
	return &Handler{
		retriever: retriever,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	localLog := zerolog.Ctx(r.Context()).With().Str("handler", "get_items").Logger()

	userID, err := converter.UserToInt(r.PathValue(constants.UserID))
	if err != nil {
		localLog.Error().Err(err).Str(constants.UserID, r.PathValue(constants.UserID)).Send()
		errhandle.NewErr(err.Error()).Send(w, localLog, http.StatusBadRequest)
		return
	}

	items, err := h.retriever.GetItemsByUserID(r.Context(), userID)
	if err != nil {
		localLog.Error().Err(err).Send()
		switch {
		case errors.Is(err, models.ErrCartIsEmpty):
			errhandle.NewErr(models.ErrCartIsEmpty.Error()).Send(w, localLog, http.StatusNotFound)
		case errors.Is(err, models.ErrNotFound):
			errhandle.NewErr(models.ErrNotFound.Error()).Send(w, localLog, http.StatusPreconditionFailed)
		default:
			errhandle.NewErr(models.ErrCartCheckout.Error()).Send(w, localLog, http.StatusInternalServerError)
		}
		return
	}

	itemsResp := itemsToDTO(items)

	data, err := json.Marshal(itemsResp)
	if err != nil {
		localLog.Error().Err(err).Send()
		errhandle.NewErr(models.ErrJSONProcessing.Error()).Send(w, localLog, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(data); err != nil {
		localLog.Error().Err(err).Send()
	}
}
