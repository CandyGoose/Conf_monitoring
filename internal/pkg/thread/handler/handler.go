package handler

import (
	"fmt"
	"server/internal/pkg/domain"
	"server/internal/pkg/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handler struct {
	ThreadSvc domain.ThreadService
	Logger    *zap.SugaredLogger
}

func (h Handler) GetThread(ctx echo.Context) error {
	h.Logger.Infow("GetThread called", utils.GetRequestData(ctx)...)

	tid := ctx.Param("tid")

	t, err := h.ThreadSvc.Get(tid)
	h.Logger.Infow(
		fmt.Sprintf("ThreadService.Get called with threadID = %s, result = %#v", tid, t),
		utils.GetRequestData(ctx)...,
	)

	if err != nil {
		h.Logger.Errorw(fmt.Sprintf("ThreadService.Get failed: %v", err), utils.GetRequestData(ctx)...)
		return err
	}

	err = ctx.JSON(200, t)
	if err != nil {
		h.Logger.Errorw(fmt.Sprintf("ctx.JSON failed: %v", err), utils.GetRequestData(ctx)...)
	}

	return err
}

func (h Handler) CreateThread(ctx echo.Context) error {
	h.Logger.Infow("CreateThread called", utils.GetRequestData(ctx)...)

	var thread domain.Thread

	err := ctx.Bind(&thread)
	if err != nil {
		h.Logger.Errorw(fmt.Sprintf("ctx.Bind failed: %v", err), utils.GetRequestData(ctx)...)
		return err
	}

	err = h.ThreadSvc.Create(thread)

	h.Logger.Infow(
		fmt.Sprintf("ThreadService.Create called with thread = %#v", thread),
		utils.GetRequestData(ctx)...,
	)

	if err != nil {
		h.Logger.Errorw(fmt.Sprintf("ThreadService.Create failed: %v", err), utils.GetRequestData(ctx)...)
		return err
	}

	return ctx.NoContent(200)
}
