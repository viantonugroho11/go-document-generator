package apis

import (
	"github.com/labstack/echo/v4"

	"go-document-generator/internal/transport/apis/handler"
	"go-document-generator/internal/usecase/documenttemplateversions"
	"go-document-generator/internal/usecase/documenttemplates"
	"go-document-generator/internal/usecase/users"
)

func RegisterRoutes(
	e *echo.Echo,
	userService users.UserService,
	tmplService documenttemplates.TemplatesService,
	verService documenttemplateversions.VersionsService,
) {
	userHandler := handler.NewUserHandler(userService)
	templateHandler := handler.NewTemplateHandler(tmplService)
	versionHandler := handler.NewTemplateVersionHandler(verService)

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(200, "ok")
	})

	users := e.Group("/users")
	users.POST("", userHandler.Create)
	users.GET("", userHandler.List)
	users.GET("/:id", userHandler.GetByID)
	users.PUT("/:id", userHandler.Update)
	users.DELETE("/:id", userHandler.Delete)

	// Templates
	templates := e.Group("/templates")
	templates.POST("", templateHandler.CreateOrAdd)
	templates.GET("", templateHandler.List)
	templates.GET("/:id", templateHandler.GetByID)
	templates.PUT("/:id", templateHandler.Update)
	templates.DELETE("/:id", templateHandler.Delete)
	templates.GET("/:id/versions", versionHandler.ListByTemplateID)
	templates.POST("/:id/versions", versionHandler.Create)

	// Template Versions
	templateVersions := e.Group("/template-versions")
	templateVersions.GET("/:versionId", versionHandler.GetByID)
	templateVersions.PUT("/:versionId", versionHandler.Update)
	templateVersions.DELETE("/:versionId", versionHandler.Delete)
}
