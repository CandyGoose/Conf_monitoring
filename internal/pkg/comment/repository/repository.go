package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/internal/metrics"
	"server/internal/pkg/domain"
	"server/internal/pkg/utils"

	"go.uber.org/zap"
)

type repository struct {
	Logger  *zap.SugaredLogger
	Metrics *metrics.HTTPRequestMetrics
}

func NewRepository(logger *zap.SugaredLogger, metrics *metrics.HTTPRequestMetrics) domain.CommentRepository {
	return repository{
		Logger:  logger,
		Metrics: metrics,
	}
}

func (r repository) Create(comment domain.Comment) error {
	reqBody, err := json.Marshal(comment)
	if err != nil {
		return fmt.Errorf("failed to marshal variable 'comment': %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://95.163.251.187:16000/comment?fast=true", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("http.NewRequest failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, latency, err := utils.SendHTTPRequest(req, r.Metrics)
	if err != nil {
		return fmt.Errorf("failed to send http request: %w", err)
	}
	defer resp.Body.Close()

	r.Logger.Infow(
		"Request to an external comment service",
		utils.GetServiceRequestData(req, resp, latency)...,
	)

	if resp.StatusCode != 200 {
		return errors.New("failed to create comment remotely")
	}

	return nil
}

func (r repository) Like(commentID string) error {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://95.163.251.187:16000/comment/like?cid=%s&superstable=true", commentID),
		nil,
	)
	if err != nil {
		return fmt.Errorf("http.NewRequest failed: %w", err)
	}

	resp, latency, err := utils.SendHTTPRequest(req, r.Metrics)
	if err != nil {
		return fmt.Errorf("failed to send http request: %w", err)
	}
	defer resp.Body.Close()

	r.Logger.Infow(
		"Request to an external comment service",
		utils.GetServiceRequestData(req, resp, latency)...,
	)

	if resp.StatusCode != 200 {
		return errors.New("failed to like comment remotely")
	}

	return nil
}
