package utils

import (
	"net/http"
	"server/internal/metrics"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func SendHTTPRequest(req *http.Request, m *metrics.HTTPRequestMetrics) (*http.Response, time.Duration, error) {
	url := req.URL
	timer := prometheus.NewTimer(m.HTTPDuration.WithLabelValues(url.Host, req.Method, url.Path))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	d := timer.ObserveDuration()
	m.HTTPRequests.WithLabelValues(url.Host, strconv.Itoa(resp.StatusCode), req.Method, url.Path).Inc()
	return resp, d, err
}
