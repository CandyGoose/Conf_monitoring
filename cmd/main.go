package main

import (
	"fmt"
	"log"
	"server/internal/api/middleware"
	"server/internal/metrics"
	"server/internal/pkg/comment/handler"
	commentrepo "server/internal/pkg/comment/repository"
	commentsvc "server/internal/pkg/comment/service"
	"server/internal/pkg/session"
	threadhttp "server/internal/pkg/thread/handler"
	threadrepo "server/internal/pkg/thread/repository"
	threadsvc "server/internal/pkg/thread/service"

	echoPrometheus "github.com/globocom/echo-prometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	zapConfig := zap.NewProductionConfig()
	zapConfig.OutputPaths = []string{"server.log"}
	zapConfig.Level.SetLevel(zap.ErrorLevel)
	zapLogger, err := zapConfig.Build()
	if err != nil {
		log.Fatalf("failed to build zap logger: %v", err)
	}
	defer func() {
		if err := zapLogger.Sync(); err != nil {
			log.Fatalf("zap logger sync failed: %v", err)
		}
	}()
	logger := zapLogger.Sugar()

	e := echo.New()

	serverConfigMetrics := metrics.NewConfig()
	serverConfigMetrics.Subsystem = "server"
	e.Use(echoPrometheus.MetricsMiddlewareWithConfig(serverConfigMetrics))

	e.Use(middleware.RequestID())
	e.Use(middleware.Recover(logger))
	e.Use(middleware.AccessLog(logger))

	serviceConfigMetrics := metrics.NewConfig()
	serviceConfigMetrics.Subsystem = "services"
	serviceMetrics := metrics.NewServiceMetrics(serviceConfigMetrics)

	threadRepo := threadrepo.NewRepository(logger, serviceMetrics)
	threadSvc := threadsvc.NewService(threadRepo)
	threadHandler := threadhttp.Handler{ThreadSvc: threadSvc, Logger: logger}

	commentRepo := commentrepo.NewRepository(logger, serviceMetrics)
	commentSvc := commentsvc.NewService(commentRepo, threadRepo)
	commentHandler := handler.Handler{CommentSvc: commentSvc, Logger: logger}

	sessionSvc := session.NewService(logger, serviceMetrics)
	auth := middleware.AuthEchoMiddleware(sessionSvc, logger)

	e.GET("/thread/:tid", threadHandler.GetThread, auth)
	e.POST("/thread", threadHandler.CreateThread, auth)
	e.POST("/thread/:tid/comment", commentHandler.Create, auth)
	e.POST("/thread/:tid/comment/:cid/like", commentHandler.Like, auth)

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	fmt.Print(e.Start(":8000"))
}
