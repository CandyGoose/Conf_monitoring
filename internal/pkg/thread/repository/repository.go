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

func NewRepository(logger *zap.SugaredLogger, metrics *metrics.HTTPRequestMetrics) domain.ThreadRepository {
	return repository{
		Logger:  logger,
		Metrics: metrics,
	}
}

func (r repository) Create(thread domain.Thread) error {
	reqBody, err := json.Marshal(thread)
	if err != nil {
		return fmt.Errorf("failed to marshal variable 'thread': %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://95.163.251.187:15000/thread?stable=true", bytes.NewBuffer(reqBody))
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
		"Request to an external thread service",
		utils.GetServiceRequestData(req, resp, latency)...,
	)

	if resp.StatusCode != 200 {
		return errors.New("failed to create thread remotely")
	}

	return nil
}

func (r repository) Get(id string) (domain.Thread, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://95.163.251.187:15000/thread?id=%s&thread_fast=true", id), nil)
	if err != nil {
		return domain.Thread{}, fmt.Errorf("http.NewRequest failed: %w", err)
	}

	resp, latency, err := utils.SendHTTPRequest(req, r.Metrics)
	if err != nil {
		return domain.Thread{}, fmt.Errorf("failed to send http request: %w", err)
	}
	defer resp.Body.Close()

	r.Logger.Infow(
		"Request to an external thread service",
		utils.GetServiceRequestData(req, resp, latency)...,
	)

	if resp.StatusCode != 200 {
		return domain.Thread{}, errors.New("failed to fetch thread remotely")
	}

	var thread domain.Thread
	err = json.NewDecoder(resp.Body).Decode(&thread)
	if err != nil {
		return domain.Thread{}, fmt.Errorf("failed to decode response body into json: %w", err)
	}

	return thread, nil
}
