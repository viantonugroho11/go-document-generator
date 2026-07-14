package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go-document-generator/internal/shared/apperror"
)

// APIKeyAuth middleware validasi API key dari header Authorization atau X-API-Key.
// Jika validKeys kosong, auth dinonaktifkan (dev mode).
func APIKeyAuth(validKeys []string) echo.MiddlewareFunc {
	if len(validKeys) == 0 {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	keySet := make(map[string]struct{}, len(validKeys))
	for _, k := range validKeys {
		keySet[strings.TrimSpace(k)] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := extractKey(c.Request())
			if key == "" {
				return c.JSON(http.StatusUnauthorized, apperror.New("UNAUTHORIZED", "missing api key"))
			}
			if _, ok := keySet[key]; !ok {
				return c.JSON(http.StatusUnauthorized, apperror.New("UNAUTHORIZED", "invalid api key"))
			}
			return next(c)
		}
	}
}

func extractKey(r *http.Request) string {
	if v := r.Header.Get("X-API-Key"); v != "" {
		return strings.TrimSpace(v)
	}
	if v := r.Header.Get("Authorization"); strings.HasPrefix(v, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(v, "Bearer "))
	}
	return ""
}
