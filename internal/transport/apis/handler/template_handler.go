package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go-document-generator/internal/entity/enums"
	tplrepo "go-document-generator/internal/repository/documenttemplates"
	"go-document-generator/internal/shared/pagination"
	"go-document-generator/internal/shared/tenant"
	"go-document-generator/internal/transport/apis/dto"
	ucTpl "go-document-generator/internal/usecase/documenttemplates"
)

type TemplateHandler struct {
	svc ucTpl.Service
}

func NewTemplateHandler(svc ucTpl.Service) *TemplateHandler {
	return &TemplateHandler{svc: svc}
}

func (h *TemplateHandler) List(c echo.Context) error {
	headerTenant, err := tenant.FromEcho(c)
	if err != nil {
		return writeError(c, err)
	}
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	var isActive *bool
	if v := c.QueryParam("is_active"); v != "" {
		b := v == "true"
		isActive = &b
	}
	f := tplrepo.ListFilter{
		TenantID: headerTenant,
		Code:     c.QueryParam("code"),
		Category: c.QueryParam("category"),
		IsActive: isActive,
		Engine:   enums.TemplateEngine(c.QueryParam("engine")),
		Page:     pagination.Params{Page: page, Limit: limit, Sort: c.QueryParam("sort")},
	}
	items, meta, err := h.svc.List(c.Request().Context(), f)
	if err != nil {
		return writeError(c, err)
	}
	data := make([]dto.DocumentTemplateResponse, len(items))
	for i, t := range items {
		data[i] = dto.TemplateFromEntity(t)
	}
	return c.JSON(http.StatusOK, dto.TemplateListResponse{Data: data, Meta: dto.MetaFrom(meta)})
}

func (h *TemplateHandler) Create(c echo.Context) error {
	headerTenant, err := tenant.FromEcho(c)
	if err != nil {
		return writeError(c, err)
	}
	var req dto.CreateTemplateRequest
	if err := c.Bind(&req); err != nil {
		return writeError(c, err)
	}
	created, err := h.svc.Create(c.Request().Context(), req.ToEntity(headerTenant))
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusCreated, dto.TemplateFromEntity(created))
}

func (h *TemplateHandler) Get(c echo.Context) error {
	headerTenant, err := tenant.FromEcho(c)
	if err != nil {
		return writeError(c, err)
	}
	id, err := strconv.ParseInt(c.Param("template_id"), 10, 64)
	if err != nil {
		return writeError(c, err)
	}
	t, err := h.svc.GetByID(c.Request().Context(), id, headerTenant)
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, dto.TemplateFromEntity(t))
}

func (h *TemplateHandler) Patch(c echo.Context) error {
	headerTenant, err := tenant.FromEcho(c)
	if err != nil {
		return writeError(c, err)
	}
	id, err := strconv.ParseInt(c.Param("template_id"), 10, 64)
	if err != nil {
		return writeError(c, err)
	}
	existing, err := h.svc.GetByID(c.Request().Context(), id, headerTenant)
	if err != nil {
		return writeError(c, err)
	}
	var req dto.PatchTemplateRequest
	if err := c.Bind(&req); err != nil {
		return writeError(c, err)
	}
	updated, err := h.svc.Patch(c.Request().Context(), dto.ApplyPatchTemplate(existing, req))
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, dto.TemplateFromEntity(updated))
}

func (h *TemplateHandler) Delete(c echo.Context) error {
	headerTenant, err := tenant.FromEcho(c)
	if err != nil {
		return writeError(c, err)
	}
	id, err := strconv.ParseInt(c.Param("template_id"), 10, 64)
	if err != nil {
		return writeError(c, err)
	}
	if err := h.svc.Deactivate(c.Request().Context(), id, headerTenant, nil); err != nil {
		return writeError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}
