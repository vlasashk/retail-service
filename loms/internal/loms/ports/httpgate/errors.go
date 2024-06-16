package httpgate

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errResp struct {
	Error string `json:"error"`
}

func errHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	s, ok := status.FromError(err)
	if !ok {
		runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
		return
	}

	switch s.Code() {
	case codes.InvalidArgument:
		customErr(w, http.StatusBadRequest, s.Message())
	case codes.NotFound:
		customErr(w, http.StatusNotFound, s.Message())
	case codes.FailedPrecondition:
		customErr(w, http.StatusPreconditionFailed, s.Message())
	case codes.Internal:
		customErr(w, http.StatusInternalServerError, s.Message())
	default:
		runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
	}

}

func customErr(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errMsg := errResp{
		Error: msg,
	}
	if err := json.NewEncoder(w).Encode(errMsg); err != nil {
		log.Error().Err(err).Send()
	}
}
