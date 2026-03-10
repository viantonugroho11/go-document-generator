package handler

import (
	"net/http"
	"strconv"

	"go-document-generator/internal/transport/apis/dto"
	useVer "go-document-generator/internal/usecase/documenttemplateversions"

	"github.com/labstack/echo/v4"
)

type TemplateVersionHandler struct {
	service useVer.VersionsService
}

func NewTemplateVersionHandler(service useVer.VersionsService) *TemplateVersionHandler {
	return &TemplateVersionHandler{service: service}
}

// POST /templates/:id/versions
func (h *TemplateVersionHandler) Create(c echo.Context) error {
	idStr := c.Param("id")
	templateID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid template id"})
	}
	var req dto.CreateTemplateVersionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	ver := req.ToEntity()
	ver.TemplateID = templateID
	created, err := h.service.Create(c.Request().Context(), ver)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, created)
}

// GET /templates/:id/versions
func (h *TemplateVersionHandler) ListByTemplateID(c echo.Context) error {
	idStr := c.Param("id")
	templateID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid template id"})
	}
	items, err := h.service.ListByTemplateID(c.Request().Context(), templateID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, items)
}

// GET /template-versions/:versionId
func (h *TemplateVersionHandler) GetByID(c echo.Context) error {
	idStr := c.Param("versionId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid version id"})
	}
	item, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, item)
}

// PUT /template-versions/:versionId
func (h *TemplateVersionHandler) Update(c echo.Context) error {
	idStr := c.Param("versionId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid version id"})
	}
	var req dto.UpdateTemplateVersionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	existing, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	existing.Content = req.Content
	existing.Schema = req.Schema
	existing.SamplePayload = req.SamplePayload
	if req.IsPublished != nil {
		existing.IsPublished = *req.IsPublished
	}
	updated, err := h.service.Update(c.Request().Context(), existing)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, updated)
}

// DELETE /template-versions/:versionId
func (h *TemplateVersionHandler) Delete(c echo.Context) error {
	idStr := c.Param("versionId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid version id"})
	}
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

