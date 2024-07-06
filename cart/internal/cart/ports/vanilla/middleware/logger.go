package middleware

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"time"

	"route256/cart/pkg/utils"

	"github.com/rs/zerolog"
)

func LoggingMiddleware(log zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var requestBody []byte
			if r.Body != nil {
				bodyBytes, err := io.ReadAll(r.Body)
				if err == nil {
					requestBody = bodyBytes
					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}
			}
			traceID, spanID := utils.ExtractTraceInfo(r.Context())
			logger := log.With().Str("trace_id", traceID).Str("span_id", spanID).Logger()

			logger.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("query", r.URL.RawQuery).
				Bytes("body", requestBody).
				Any("headers", r.Header).
				Send()

			ww := &statusWriter{
				statusCode: http.StatusOK,
				err:        bytes.NewBuffer(nil),
				w:          w,
			}

			defer func(start time.Time) {
				logResp := logger.With().
					Int("status_code", ww.statusCode).
					Dur("duration", time.Since(start)).
					Any("headers", ww.Header()).
					Logger()
				switch {
				case ww.statusCode >= 400 && ww.statusCode < 500:
					logResp.Warn().Err(errors.New(ww.err.String())).Send()
				case ww.statusCode >= 500:
					logResp.Error().Err(errors.New(ww.err.String())).Send()
				default:
					logResp.Info().Send()
				}
			}(time.Now())

			next.ServeHTTP(ww, r.WithContext(logger.WithContext(r.Context())))
		})
	}
}
