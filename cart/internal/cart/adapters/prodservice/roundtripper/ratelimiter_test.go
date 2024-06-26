package roundtripper

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

type MockRoundTripper struct {
	response *http.Response
	err      error
}

func (m *MockRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func TestRateLimiter(t *testing.T) {
	logger := zerolog.Nop()

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       http.NoBody,
	}

	mockRoundTripper := &MockRoundTripper{
		response: mockResponse,
	}

	rateLimitedRoundTripper := Limit(logger, rate.Limit(10), 1)(mockRoundTripper)

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	start := time.Now()

	wg := sync.WaitGroup{}

	wg.Add(15)
	for i := 0; i < 15; i++ {
		go func() {
			resp, err := rateLimitedRoundTripper.RoundTrip(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			wg.Done()
		}()
	}
	wg.Wait()

	elapsed := time.Since(start)

	assert.GreaterOrEqual(t, elapsed, time.Second, "Rate limiter did not delay requests properly")
}
