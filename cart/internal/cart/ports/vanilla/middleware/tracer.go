package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel"
)

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer("").Start(r.Context(), r.Method+" "+r.URL.Path)
		defer span.End()

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
