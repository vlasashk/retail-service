package errhandle

import (
	"encoding/json"
	"net/http"

	"route256/cart/internal/cart/models"

	"github.com/rs/zerolog"
)

type ErrResp struct {
	Error string `json:"error"`
}

func NewErr(err string) ErrResp {
	return ErrResp{
		Error: err,
	}
}

func (resp ErrResp) Send(w http.ResponseWriter, log zerolog.Logger, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error().Err(err).Msg(models.ErrJSONProcessing.Error())
	}
}
