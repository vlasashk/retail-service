package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"status", "url"},
	)
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status", "url"},
	)

	ExternalRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "external_requests_total",
			Help: "Total number of requests to external resources",
		},
		[]string{"service", "status", "url"},
	)

	ExternalRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "external_request_duration_seconds",
			Help:    "Duration of requests to external resources",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "status", "url"},
	)
)
