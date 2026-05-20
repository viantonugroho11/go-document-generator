package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go-document-generator/internal/shared/tenant"
	"go-document-generator/internal/transport/apis/dto"
	ucVer "go-document-generator/internal/usecase/documenttemplateversions"
)

type TemplateVersionHandler struct {
	svc ucVer.Service
}

func NewTemplateVersionHandler(svc ucVer.Service) *TemplateVersionHandler {
	return &TemplateVersionHandler{svc: svc}
}

func (h *TemplateVersionHandler) List(c echo.Context) error {
	headerTenant, err := tenant.FromEcho(c)
	if err != nil {
		return writeError(c, err)
	}
	templateID, err := strconv.ParseInt(c.Param("template_id"), 10, 64)
	if err != nil {
		return writeError(c, err)
	}
	var isPublished *bool
	if v := c.QueryParam("is_published"); v != "" {
		b := v == "true"
		isPublished = &b
	}
	items, err := h.svc.List(c.Request().Context(), templateID, headerTenant, isPublished)
	if err != nil {
		return writeError(c, err)
	}
	data := make([]dto.TemplateVersionResponse, len(items))
	for i, v := range items {
		data[i] = dto.VersionFromEntity(v, false)
	}
	return c.JSON(http.StatusOK, dto.TemplateVersionListResponse{Data: data})
}

func (h *TemplateVersionHandler) Create(c echo.Context) error {
	headerTenant, err := tenant.FromEcho(c)
	if err != nil {
		return writeError(c, err)
	}
	templateID, err := strconv.ParseInt(c.Param("template_id"), 10, 64)
	if err != nil {
		return writeError(c, err)
	}
	var req dto.CreateTemplateVersionRequest
	if err := c.Bind(&req); err != nil {
		return writeError(c, err)
	}
	created, err := h.svc.Create(c.Request().Context(), templateID, headerTenant, req.ToEntity(headerTenant, templateID))
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusCreated, dto.VersionFromEntity(created, true))
}

func (h *TemplateVersionHandler) Get(c echo.Context) error {
	headerTenant, err := tenant.FromEcho(c)
	if err != nil {
		return writeError(c, err)
	}
	templateID, _ := strconv.ParseInt(c.Param("template_id"), 10, 64)
	versionID, _ := strconv.ParseInt(c.Param("version_id"), 10, 64)
	v, err := h.svc.GetByID(c.Request().Context(), templateID, versionID, headerTenant)
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, dto.VersionFromEntity(v, true))
}

func (h *TemplateVersionHandler) Publish(c echo.Context) error {
	headerTenant, err := tenant.FromEcho(c)
	if err != nil {
		return writeError(c, err)
	}
	templateID, _ := strconv.ParseInt(c.Param("template_id"), 10, 64)
	versionID, _ := strconv.ParseInt(c.Param("version_id"), 10, 64)
	v, err := h.svc.Publish(c.Request().Context(), templateID, versionID, headerTenant)
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, dto.VersionFromEntity(v, true))
}
