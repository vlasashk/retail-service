package middleware

import (
	"net/http"
	"runtime/debug"

	"route256/cart/internal/cart/constants"
	"route256/cart/internal/cart/ports/vanilla/handlers/errhandle"

	"github.com/rs/zerolog"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := *zerolog.Ctx(r.Context())

		defer func() {
			if rec := recover(); rec != nil {
				//nolint:errorlint // rec is not an error type
				if rec == http.ErrAbortHandler {
					panic(rec)
				}
				logger.Error().Any("panic", rec).Bytes("stack", debug.Stack()).Send()
				errhandle.NewErr(constants.ErrInternalError).Send(w, logger, http.StatusInternalServerError)
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}
