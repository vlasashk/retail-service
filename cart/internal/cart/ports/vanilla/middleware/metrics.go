package middleware

import (
	"bytes"
	"net/http"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"route256/cart/internal/cart/metrics"
)

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := &statusWriter{
			statusCode: http.StatusOK,
			err:        bytes.NewBuffer(nil),
			w:          w,
		}

		defer func(start time.Time) {
			duration := time.Since(start).Seconds()
			if r == nil {
				return
			}

			pattern := r.Method + " " + r.URL.Path

			// Workaround что бы получить паттерн ендпойнта (r.pat.str)
			pat := reflect.ValueOf(r).Elem().FieldByName("pat")
			// Обработка кейса когда у ендпоинт есть патер
			if !pat.IsNil() {
				patData := pat.Elem().FieldByName("str")
				pattern = reflect.NewAt(patData.Type(), unsafe.Pointer(patData.UnsafeAddr())).Elem().String()
			}

			metrics.RequestDuration.WithLabelValues(strconv.Itoa(ww.statusCode), pattern).Observe(duration)
			metrics.RequestsTotal.WithLabelValues(strconv.Itoa(ww.statusCode), pattern).Inc()

		}(time.Now())

		next.ServeHTTP(ww, r)
	})
}
