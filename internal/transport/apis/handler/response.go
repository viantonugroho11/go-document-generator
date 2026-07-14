package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go-document-generator/internal/shared/apperror"
)


func writeError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, apperror.ErrNotFound):
		return c.JSON(http.StatusNotFound, apperror.New("NOT_FOUND", err.Error()))
	case errors.Is(err, apperror.ErrConflict):
		return c.JSON(http.StatusConflict, apperror.New("CONFLICT", err.Error()))
	case errors.Is(err, apperror.ErrInvalidState):
		return c.JSON(http.StatusConflict, apperror.New("INVALID_STATE", err.Error()))
	case errors.Is(err, apperror.ErrInvalidInput):
		return c.JSON(http.StatusBadRequest, apperror.New("BAD_REQUEST", err.Error()))
	default:
		return c.JSON(http.StatusBadRequest, apperror.New("BAD_REQUEST", err.Error()))
	}
}
