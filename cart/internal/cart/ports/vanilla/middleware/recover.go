package middleware

import (
	"net/http"
	"runtime/debug"

	"route256/cart/internal/cart/ports/vanilla/handlers/errhandle"
	"route256/cart/pkg/constants"

	"github.com/rs/zerolog"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := *zerolog.Ctx(r.Context())

		defer func() {
			if err := recover(); err != nil {
				if err == http.ErrAbortHandler {
					panic(err)
				}
				logger.Error().Any("panic", err).Bytes("stack", debug.Stack()).Send()
				errhandle.NewErr(constants.ErrInternalError).Send(w, logger, http.StatusInternalServerError)
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}
