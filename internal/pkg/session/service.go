package session

import (
	"errors"
	"fmt"
	"net/http"
	"server/internal/metrics"
	"server/internal/pkg/domain"
	"server/internal/pkg/utils"

	"go.uber.org/zap"
)

type service struct {
	Logger  *zap.SugaredLogger
	Metrics *metrics.HTTPRequestMetrics
}

func NewService(logger *zap.SugaredLogger, metrics *metrics.HTTPRequestMetrics) domain.SessionService {
	return service{
		Logger:  logger,
		Metrics: metrics,
	}
}

func (s service) CheckSession(headers http.Header) (domain.Session, error) {
	req, err := http.NewRequest(http.MethodGet, "http://95.163.251.187:17000/int/CheckSession?stable=true&light=true", nil)
	if err != nil {
		return domain.Session{}, fmt.Errorf("http.NewRequest failed: %w", err)
	}

	req.Header = headers

	resp, latency, err := utils.SendHTTPRequest(req, s.Metrics)
	if err != nil {
		return domain.Session{}, fmt.Errorf("failed to send http request: %w", err)
	}
	defer resp.Body.Close()

	s.Logger.Infow(
		"Request to an external session service",
		utils.GetServiceRequestData(req, resp, latency)...,
	)

	switch resp.StatusCode {
	case 500:
		return domain.Session{}, errors.New("failed to request check session")
	case 200:
		return domain.Session{}, nil
	default:
		return domain.Session{}, domain.ErrNoSession
	}
}
