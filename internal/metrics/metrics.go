package metrics

import (
	echoPrometheus "github.com/globocom/echo-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type HTTPRequestMetrics struct {
	HTTPRequests *prometheus.CounterVec
	HTTPDuration *prometheus.HistogramVec
}

func NewServiceMetrics(config echoPrometheus.Config) *HTTPRequestMetrics {
	httpRequests := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: config.Namespace,
		Subsystem: config.Subsystem,
		Name:      "requests_total",
		Help:      "Number of http requests to services",
	}, []string{"host", "status", "method", "handler"})

	httpDuration := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: config.Namespace,
		Subsystem: config.Subsystem,
		Name:      "request_duration_seconds",
		Help:      "Time spent accessing services",
		Buckets:   config.Buckets,
	}, []string{"host", "method", "handler"})

	return &HTTPRequestMetrics{
		HTTPRequests: httpRequests,
		HTTPDuration: httpDuration,
	}
}

func NewConfig() echoPrometheus.Config {
	return echoPrometheus.Config{
		Namespace: "http",
		Buckets: []float64{
			0.005, // 5ms
			0.01,  // 10ms
			0.02,  // 20ms
			0.05,  // 50ms
			0.1,   // 100ms
			0.2,   // 200ms
			0.5,   // 500ms
			1.0,   // 1s
			2.0,   // 2s
			5.0,   // 5s
			10.0,  // 10s
			30.0,  // 30s
		},
		NormalizeHTTPStatus: false,
	}
}
