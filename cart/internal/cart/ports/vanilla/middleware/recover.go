package middleware

import (
	"net/http"

	"route256/cart/internal/cart/ports/vanilla/handlers/errhandle"
	"route256/cart/pkg/constants"

	"github.com/rs/zerolog"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := *zerolog.Ctx(r.Context())

		defer func() {
			if isPanic := recover(); isPanic != nil {
				errhandle.NewErr(constants.ErrInternalError).Send(w, logger, http.StatusInternalServerError)
				if isPanic == http.ErrAbortHandler {
					panic(isPanic)
				}
				logger.Panic().Any("panic", isPanic).Send()

			}
		}()

		next.ServeHTTP(w, r)
	})
}
