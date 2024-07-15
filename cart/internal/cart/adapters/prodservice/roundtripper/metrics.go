package roundtripper

import (
	"net/http"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"route256/cart/internal/cart/metrics"
)

const serviceName = "product-service"

type metric struct {
	next http.RoundTripper
}

func Metrics() func(next http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return metric{
			next: next,
		}
	}
}
func (m metric) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	defer func(start time.Time) {
		duration := time.Since(start).Seconds()
		if req == nil {
			return
		}

		pattern := req.Method + " " + req.URL.Path

		// Workaround что бы получить паттерн ендпойнта (r.pat.str)
		pat := reflect.ValueOf(req).Elem().FieldByName("pat")
		// Обработка кейса когда у ендпоинта есть паттерн
		if !pat.IsNil() {
			patData := pat.Elem().FieldByName("str")
			pattern = reflect.NewAt(patData.Type(), unsafe.Pointer(patData.UnsafeAddr())).Elem().String()
		}

		metrics.ExternalRequestDuration.WithLabelValues(serviceName, strconv.Itoa(resp.StatusCode), pattern).Observe(duration)
		metrics.ExternalRequestsTotal.WithLabelValues(serviceName, strconv.Itoa(resp.StatusCode), pattern).Inc()

	}(time.Now())

	resp, err = m.next.RoundTrip(req)

	return resp, err
}
