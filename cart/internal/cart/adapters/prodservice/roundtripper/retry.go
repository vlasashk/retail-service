package roundtripper

import (
	"net/http"
	"time"
)

type retryer struct {
	retries int
	next    http.RoundTripper
}

func Retry(retries int) func(next http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return retryer{
			retries: retries,
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
			return nil, err
		}
		if resp.StatusCode != 420 && resp.StatusCode != 429 {
			break
		}
		if sleep <= 3 {
			sleep++
		}
		time.Sleep(time.Duration(sleep) * time.Second)
	}

	return resp, err
}
