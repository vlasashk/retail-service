package clientbuilder

import (
	"net/http"
	"time"
)

type RoundTripperMiddleware func(http.RoundTripper) http.RoundTripper

type ClientBuilder struct {
	client        *http.Client
	roundTrippers []RoundTripperMiddleware
}

func New(timeOut time.Duration) *ClientBuilder {
	return &ClientBuilder{
		client: &http.Client{
			Timeout: timeOut,
		},
	}
}

func (cb *ClientBuilder) Use(m ...RoundTripperMiddleware) *ClientBuilder {
	if cb.roundTrippers == nil {
		cb.roundTrippers = make([]RoundTripperMiddleware, 0, len(m))
	}
	cb.roundTrippers = append(cb.roundTrippers, m...)

	return cb
}

func (cb *ClientBuilder) Build() *http.Client {
	roundTripper := chain(http.DefaultTransport, cb.roundTrippers...)
	cb.client.Transport = roundTripper
	return cb.client
}

// chain оборачивает middleware в обратном порядке, так что первый в списке оборачивает все последующие
func chain(h http.RoundTripper, m ...RoundTripperMiddleware) http.RoundTripper {
	if len(m) < 1 {
		return h
	}
	newHandler := h
	for i := len(m) - 1; i >= 0; i-- {
		newHandler = m[i](newHandler)
	}
	return newHandler
}
