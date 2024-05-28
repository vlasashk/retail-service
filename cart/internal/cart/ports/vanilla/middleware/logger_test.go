package middleware_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"route256/cart/internal/cart/ports/vanilla/middleware"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func dummyOkHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hello, world!"))
	if err != nil {
		log.Fatalf("dummyHandler error :%v", err)
	}
}

func dummyWarnHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte("bad request"))
	if err != nil {
		log.Fatalf("dummyHandler error :%v", err)
	}
}

func dummyErrHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte("internal error"))
	if err != nil {
		log.Fatalf("dummyHandler error :%v", err)
	}
}

func TestServerLoggingOK(t *testing.T) {
	testCases := []struct {
		name           string
		expectedStatus int
		handler        http.Handler
		body           string
		asserts        func(string2 string)
	}{
		{
			name:           "Ok",
			expectedStatus: http.StatusOK,
			handler:        http.HandlerFunc(dummyOkHandler),
			body:           "Hello, world!",
			asserts: func(logOutput string) {
				assert.Contains(t, logOutput, `"status_code":200`)
				assert.Contains(t, logOutput, `"level":"info"`)
				assert.NotContains(t, logOutput, `"error":`)
			},
		},
		{
			name:           "Warn",
			expectedStatus: http.StatusBadRequest,
			handler:        http.HandlerFunc(dummyWarnHandler),
			body:           "bad request",
			asserts: func(logOutput string) {
				assert.Contains(t, logOutput, `"status_code":400`)
				assert.Contains(t, logOutput, `"level":"warn"`)
				assert.Contains(t, logOutput, `"error":"bad request"`)
			},
		},
		{
			name:           "Error",
			expectedStatus: http.StatusInternalServerError,
			handler:        http.HandlerFunc(dummyErrHandler),
			body:           "internal error",
			asserts: func(logOutput string) {
				assert.Contains(t, logOutput, `"status_code":500`)
				assert.Contains(t, logOutput, `"level":"error"`)
				assert.Contains(t, logOutput, `"error":"internal error"`)
			},
		},
	}

	var buf bytes.Buffer
	testLogger := zerolog.New(&buf).With().Timestamp().Logger()

	for _, test := range testCases {
		req, err := http.NewRequest("GET", "/test?value=hello", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("User-Agent", "Test")
		req.Header.Set("Accept", "text/plain")

		ww := httptest.NewRecorder()
		wrappedHandler := middleware.LoggingMiddleware(testLogger)(test.handler)
		wrappedHandler.ServeHTTP(ww, req)

		assert.Equal(t, test.expectedStatus, ww.Code)
		assert.Equal(t, test.body, ww.Body.String())
		assert.Equal(t, "text/plain", ww.Header().Get("Content-Type"))

		logOutput := buf.String()
		println(logOutput)
		assert.Contains(t, logOutput, `"method":"GET"`)
		assert.Contains(t, logOutput, `"path":"/test"`)
		assert.Contains(t, logOutput, `"query":"value=hello"`)
		test.asserts(logOutput)
		buf.Reset()
	}
}
