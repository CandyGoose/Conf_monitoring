package handler

import (
	"fmt"
	"server/internal/pkg/domain"
	"server/internal/pkg/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handler struct {
	CommentSvc domain.CommentService
	Logger     *zap.SugaredLogger
}

func (h Handler) Create(ctx echo.Context) error {
	h.Logger.Infow("Create called", utils.GetRequestData(ctx)...)

	var comment domain.Comment

	err := ctx.Bind(&comment)
	if err != nil {
		h.Logger.Errorw(fmt.Sprintf("ctx.Bind failed: %v", err), utils.GetRequestData(ctx)...)
		return err
	}

	tid := ctx.Param("tid")
	err = h.CommentSvc.Create(tid, comment)

	h.Logger.Infow(
		fmt.Sprintf("CommentService.Create called with threadID = %s, comment = %#v", tid, comment),
		utils.GetRequestData(ctx)...,
	)

	if err != nil {
		h.Logger.Errorw(
			fmt.Sprintf("CommentService.Create failed: %v", err),
			utils.GetRequestData(ctx)...,
		)
		return err
	}

	return ctx.NoContent(200)
}

func (h Handler) Like(ctx echo.Context) error {
	h.Logger.Infow("Like called", utils.GetRequestData(ctx)...)

	tid := ctx.Param("tid")
	cid := ctx.Param("cid")

	err := h.CommentSvc.Like(tid, cid)
	logMessage := fmt.Sprintf("CommentService.Like called with threadID = %s, commentID = %s", tid, cid)

	if err != nil {
		h.Logger.Errorw(fmt.Sprintf("CommentService.Like failed: %v", err), append(utils.GetRequestData(ctx), "error", err.Error())...)
		return err
	}

	h.Logger.Infow(logMessage, utils.GetRequestData(ctx)...)
	return ctx.NoContent(200)
}
