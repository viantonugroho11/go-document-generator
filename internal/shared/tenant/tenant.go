package tenant

import (
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const HeaderTenantID = "X-Tenant-Id"

func FromEcho(c echo.Context) (*string, error) {
	raw := strings.TrimSpace(c.Request().Header.Get(HeaderTenantID))
	if raw == "" {
		return nil, nil
	}
	if _, err := uuid.Parse(raw); err != nil {
		return nil, err
	}
	return &raw, nil
}

func PtrOrNil(id string) *string {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil
	}
	return &id
}
