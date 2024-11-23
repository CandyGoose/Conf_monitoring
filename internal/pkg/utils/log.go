package utils

import (
	"net/http"
	"server/internal/api/middleware"
	"time"

	"github.com/labstack/echo/v4"
)

func GetRequestData(ctx echo.Context) []interface{} {
	return []interface{}{
		"request_id", ctx.Get(middleware.RequestIDKey),
		"remote_addr", ctx.Request().RemoteAddr,
		"method", ctx.Request().Method,
		"url", ctx.Request().URL,
	}
}

func GetServiceRequestData(req *http.Request, resp *http.Response, latency time.Duration) []interface{} {
	return []interface{}{
		"method", req.Method,
		"url", req.URL,
		"status", resp.StatusCode,
		"latency", latency.String(),
	}
}
