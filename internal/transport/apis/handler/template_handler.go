package handler

import (
	"net/http"
	"strconv"

	"go-document-generator/internal/transport/apis/dto"
	useTmpl "go-document-generator/internal/usecase/documenttemplates"

	"github.com/labstack/echo/v4"
)

type TemplateHandler struct {
	service useTmpl.TemplatesService
}

func NewTemplateHandler(service useTmpl.TemplatesService) *TemplateHandler {
	return &TemplateHandler{service: service}
}

// POST /templates
// Membuat template baru (dengan versi 1) atau menambah versi pada template existing (berdasarkan code).
func (h *TemplateHandler) CreateOrAdd(c echo.Context) error {
	var req dto.CreateOrAddTemplateWithVersionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	tmpl, ver, err := h.service.CreateOrAddVersion(c.Request().Context(), req.ToTemplateEntity(), req.ToVersionEntity())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]any{
		"template": tmpl,
		"version":  ver,
	})
}

func (h *TemplateHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	tmpl, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, tmpl)
}

func (h *TemplateHandler) List(c echo.Context) error {
	items, err := h.service.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, items)
}

func (h *TemplateHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	var req dto.UpdateTemplateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	tmpl, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	// apply updates
	tmpl.Name = req.Name
	tmpl.Description = req.Description
	tmpl.Engine = req.Engine
	tmpl.OutputFormat = req.OutputFormat
	if req.IsActive != nil {
		tmpl.IsActive = *req.IsActive
	}
	updated, err := h.service.Update(c.Request().Context(), tmpl)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, updated)
}

func (h *TemplateHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

