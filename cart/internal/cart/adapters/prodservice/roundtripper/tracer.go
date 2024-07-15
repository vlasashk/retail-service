package roundtripper

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type tracer struct {
	next http.RoundTripper
}

func Tracing() func(next http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return tracer{
			next: next,
		}
	}
}

func (t tracer) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx, span := otel.Tracer("").Start(req.Context(), req.Method+" "+req.URL.Path)
	defer span.End()

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := t.next.RoundTrip(req.WithContext(ctx))
	if err != nil {
		span.RecordError(err)
	}

	return resp, err
}
