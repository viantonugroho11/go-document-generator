package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go-document-generator/internal/transport/apis/dto"
	ucCb "go-document-generator/internal/usecase/documentcallbackattempts"
)

type CallbackHandler struct {
	svc ucCb.Service
}

func NewCallbackHandler(svc ucCb.Service) *CallbackHandler {
	return &CallbackHandler{svc: svc}
}

func (h *CallbackHandler) Test(c echo.Context) error {
	var req dto.TestCallbackRequest
	if err := c.Bind(&req); err != nil {
		return writeError(c, err)
	}
	result, err := h.svc.TestCallback(c.Request().Context(), ucCb.TestCallbackInput{
		CallbackURL: req.CallbackURL, SamplePayload: req.SamplePayload,
	})
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, dto.TestCallbackResponse{
		Success: result.Success, ResponseStatusCode: result.ResponseStatusCode,
		ErrorMessage: result.ErrorMessage,
	})
}
