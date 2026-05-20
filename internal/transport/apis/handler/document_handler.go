package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go-document-generator/internal/entity/enums"
	docrepo "go-document-generator/internal/repository/documents"
	"go-document-generator/internal/shared/pagination"
	"go-document-generator/internal/shared/tenant"
	"go-document-generator/internal/transport/apis/dto"
	ucCb "go-document-generator/internal/usecase/documentcallbackattempts"
	ucDoc "go-document-generator/internal/usecase/documents"
	ucLog "go-document-generator/internal/usecase/documentrenderlogs"
)

type DocumentHandler struct {
	docs     ucDoc.Service
	logs     ucLog.Service
	callback ucCb.Service
}

func NewDocumentHandler(docs ucDoc.Service, logs ucLog.Service, callback ucCb.Service) *DocumentHandler {
	return &DocumentHandler{docs: docs, logs: logs, callback: callback}
}

func (h *DocumentHandler) List(c echo.Context) error {
	headerTenant, err := tenant.FromEcho(c)
	if err != nil {
		return writeError(c, err)
	}
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	f := docrepo.ListFilter{
		TenantID:     headerTenant,
		RequestID:    c.QueryParam("request_id"),
		Status:       enums.DocumentStatus(c.QueryParam("status")),
		TemplateCode: c.QueryParam("template_code"),
		DmsStatus:    enums.DmsStatus(c.QueryParam("dms_status")),
		CallbackStatus: enums.CallbackStatus(c.QueryParam("callback_status")),
		Page: pagination.Params{Page: page, Limit: limit, Sort: c.QueryParam("sort")},
	}
	if v := c.QueryParam("created_from"); v != "" {
		if t, e := time.Parse(time.RFC3339, v); e == nil {
			f.CreatedFrom = &t
		}
	}
	if v := c.QueryParam("created_to"); v != "" {
		if t, e := time.Parse(time.RFC3339, v); e == nil {
			f.CreatedTo = &t
		}
	}
	items, meta, err := h.docs.List(c.Request().Context(), f)
	if err != nil {
		return writeError(c, err)
	}
	data := make([]dto.GeneratedDocumentResponse, len(items))
	for i, d := range items {
		data[i] = dto.DocumentFromEntity(d)
	}
	return c.JSON(http.StatusOK, dto.DocumentListResponse{Data: data, Meta: dto.MetaFrom(meta)})
}

func (h *DocumentHandler) Create(c echo.Context) error {
	headerTenant, err := tenant.FromEcho(c)
	if err != nil {
		return writeError(c, err)
	}
	var req dto.CreateDocumentRequest
	if err := c.Bind(&req); err != nil {
		return writeError(c, err)
	}
	doc, replay, err := h.docs.Create(c.Request().Context(), req.ToInput(headerTenant))
	if err != nil {
		return writeError(c, err)
	}
	status := http.StatusAccepted
	if replay {
		status = http.StatusOK
	}
	return c.JSON(status, dto.DocumentFromEntity(doc))
}

func (h *DocumentHandler) Get(c echo.Context) error {
	headerTenant, err := tenant.FromEcho(c)
	if err != nil {
		return writeError(c, err)
	}
	id, _ := strconv.ParseInt(c.Param("document_id"), 10, 64)
	doc, err := h.docs.GetByID(c.Request().Context(), id, headerTenant)
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, dto.DocumentFromEntity(doc))
}

func (h *DocumentHandler) Patch(c echo.Context) error {
	headerTenant, err := tenant.FromEcho(c)
	if err != nil {
		return writeError(c, err)
	}
	id, err := strconv.ParseInt(c.Param("document_id"), 10, 64)
	if err != nil {
		return writeError(c, err)
	}
	existing, err := h.docs.GetByID(c.Request().Context(), id, headerTenant)
	if err != nil {
		return writeError(c, err)
	}
	var req dto.PatchDocumentRequest
	if err := c.Bind(&req); err != nil {
		return writeError(c, err)
	}
	updated, err := h.docs.Patch(c.Request().Context(), dto.ApplyPatchDocument(existing, req))
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, dto.DocumentFromEntity(updated))
}

func (h *DocumentHandler) GetByRequestID(c echo.Context) error {
	headerTenant, err := tenant.FromEcho(c)
	if err != nil {
		return writeError(c, err)
	}
	doc, err := h.docs.GetByRequestID(c.Request().Context(), c.Param("request_id"), headerTenant)
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, dto.DocumentFromEntity(doc))
}

func (h *DocumentHandler) Delete(c echo.Context) error {
	headerTenant, err := tenant.FromEcho(c)
	if err != nil {
		return writeError(c, err)
	}
	id, _ := strconv.ParseInt(c.Param("document_id"), 10, 64)
	if err := h.docs.SoftDelete(c.Request().Context(), id, headerTenant); err != nil {
		return writeError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *DocumentHandler) Cancel(c echo.Context) error {
	headerTenant, err := tenant.FromEcho(c)
	if err != nil {
		return writeError(c, err)
	}
	id, _ := strconv.ParseInt(c.Param("document_id"), 10, 64)
	doc, err := h.docs.Cancel(c.Request().Context(), id, headerTenant)
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, dto.DocumentFromEntity(doc))
}

func (h *DocumentHandler) Retry(c echo.Context) error {
	headerTenant, err := tenant.FromEcho(c)
	if err != nil {
		return writeError(c, err)
	}
	id, _ := strconv.ParseInt(c.Param("document_id"), 10, 64)
	doc, err := h.docs.Retry(c.Request().Context(), id, headerTenant)
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusAccepted, dto.DocumentFromEntity(doc))
}

func (h *DocumentHandler) Download(c echo.Context) error {
	headerTenant, err := tenant.FromEcho(c)
	if err != nil {
		return writeError(c, err)
	}
	id, _ := strconv.ParseInt(c.Param("document_id"), 10, 64)
	url, err := h.docs.DownloadURL(c.Request().Context(), id, headerTenant)
	if err != nil {
		return writeError(c, err)
	}
	return c.Redirect(http.StatusFound, url)
}

func (h *DocumentHandler) ListRenderLogs(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("document_id"), 10, 64)
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	items, meta, err := h.logs.ListByDocumentID(c.Request().Context(), id, pagination.Params{Page: page, Limit: limit})
	if err != nil {
		return writeError(c, err)
	}
	data := make([]dto.RenderLogResponse, len(items))
	for i, l := range items {
		data[i] = dto.RenderLogFromEntity(l)
	}
	return c.JSON(http.StatusOK, dto.RenderLogListResponse{Data: data, Meta: dto.MetaFrom(meta)})
}

func (h *DocumentHandler) ListCallbackAttempts(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("document_id"), 10, 64)
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	items, meta, err := h.callback.ListByDocumentID(c.Request().Context(), id, pagination.Params{Page: page, Limit: limit})
	if err != nil {
		return writeError(c, err)
	}
	data := make([]dto.CallbackAttemptResponse, len(items))
	for i, a := range items {
		data[i] = dto.CallbackFromEntity(a)
	}
	return c.JSON(http.StatusOK, dto.CallbackAttemptListResponse{Data: data, Meta: dto.MetaFrom(meta)})
}
