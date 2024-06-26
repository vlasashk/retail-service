package roundtripper

import (
	"net/http"

	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

type limiter struct {
	rateLimiter *rate.Limiter
	log         zerolog.Logger
	next        http.RoundTripper
}

func Limit(log zerolog.Logger, rpsLimit rate.Limit, burstLimit int) func(next http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return limiter{
			rateLimiter: rate.NewLimiter(rpsLimit, burstLimit),
			log:         log,
			next:        next,
		}
	}
}

func (l limiter) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := l.rateLimiter.Wait(req.Context()); err != nil {
		return nil, err
	}

	return l.next.RoundTrip(req)
}
