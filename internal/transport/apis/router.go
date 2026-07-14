package apis

import (
	"github.com/labstack/echo/v4"

	"go-document-generator/internal/transport/apis/handler"
	ucCb "go-document-generator/internal/usecase/documentcallbackattempts"
	ucDoc "go-document-generator/internal/usecase/documents"
	ucLog "go-document-generator/internal/usecase/documentrenderlogs"
	ucTpl "go-document-generator/internal/usecase/documenttemplates"
	ucVer "go-document-generator/internal/usecase/documenttemplateversions"
	usecaseusers "go-document-generator/internal/usecase/users"
)

type Services struct {
	Users            usecaseusers.UserService
	Templates        ucTpl.Service
	TemplateVersions ucVer.Service
	Documents        ucDoc.Service
	RenderLogs       ucLog.Service
	Callbacks        ucCb.Service
}

func RegisterRoutes(e *echo.Echo, svc Services) {
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(200, "ok")
	})

	if svc.Users != nil {
		userHandler := handler.NewUserHandler(svc.Users)
		users := e.Group("/users")
		users.POST("", userHandler.Create)
		users.GET("", userHandler.List)
		users.GET("/:id", userHandler.GetByID)
		users.PUT("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
	}

	tplHandler := handler.NewTemplateHandler(svc.Templates)
	verHandler := handler.NewTemplateVersionHandler(svc.TemplateVersions)
	docHandler := handler.NewDocumentHandler(svc.Documents, svc.RenderLogs, svc.Callbacks)
	cbHandler := handler.NewCallbackHandler(svc.Callbacks)

	templates := e.Group("/templates")
	templates.GET("", tplHandler.List)
	templates.POST("", tplHandler.Create)
	templates.GET("/:template_id", tplHandler.Get)
	templates.PATCH("/:template_id", tplHandler.Patch)
	templates.DELETE("/:template_id", tplHandler.Delete)

	templates.GET("/:template_id/versions", verHandler.List)
	templates.POST("/:template_id/versions", verHandler.Create)
	templates.GET("/:template_id/versions/:version_id", verHandler.Get)
	templates.POST("/:template_id/versions/:version_id/publish", verHandler.Publish)
	templates.POST("/:template_id/versions/:version_id/preview", docHandler.Preview)

	docs := e.Group("/documents")
	docs.GET("", docHandler.List)
	docs.POST("", docHandler.Create)
	docs.POST("/bulk", docHandler.BulkCreate)
	docs.GET("/by-request/:request_id", docHandler.GetByRequestID)
	docs.GET("/:document_id", docHandler.Get)
	docs.PATCH("/:document_id", docHandler.Patch)
	docs.DELETE("/:document_id", docHandler.Delete)
	docs.POST("/:document_id/cancel", docHandler.Cancel)
	docs.POST("/:document_id/retry", docHandler.Retry)
	docs.GET("/:document_id/download", docHandler.Download)
	docs.GET("/:document_id/render-logs", docHandler.ListRenderLogs)
	docs.GET("/:document_id/callback-attempts", docHandler.ListCallbackAttempts)

	e.POST("/callbacks/test", cbHandler.Test)
}
