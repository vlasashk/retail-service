package roundtripper

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type retryer struct {
	retries int
	log     zerolog.Logger
	next    http.RoundTripper
}

func Retry(log zerolog.Logger, retries int) func(next http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return retryer{
			retries: retries,
			log:     log,
			next:    next,
		}
	}
}

func (r retryer) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)

	sleep := 0
	for i := 0; i < r.retries; i++ {
		resp, err = r.next.RoundTrip(req)
		if err != nil {
			log.Error().Err(err).Send()
			return nil, err
		}
		if resp.StatusCode != 420 && resp.StatusCode != 429 {
			break
		}
		log.Warn().Int("retries left", r.retries-1-i).Send()
		if sleep <= 3 {
			sleep++
		}
		if i != r.retries-1 {
			time.Sleep(time.Duration(sleep) * time.Second)
		}
	}

	return resp, err
}
