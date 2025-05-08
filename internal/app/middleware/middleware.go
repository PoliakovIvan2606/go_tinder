package middleware

import (
	"net/http"
	"strings"
	"time"
	"tinder/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func isSkippablePath(path string) bool {
	skipPaths := map[string]bool{
		"/login":           true,
		"/register":        true,
		"/api/user/login":  true,
		"/api/user/add":    true,
	}

	// Пропуск статических файлов
	staticPrefixes := []string{
		"/static/",
		"/assets/",
		"/js/",
		"/css/",
		"/upload/",
	}

	if skipPaths[path] {
		return true
	}

	for _, prefix := range staticPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isSkippablePath(c.Request.URL.Path) {
			c.Next()
			return
		}

		accessCookie, err := c.Request.Cookie("access")
		refreshCookie, rErr := c.Request.Cookie("refresh")

		if err == nil {
			if claims, ok := jwt.ValidateJWT(accessCookie.Value); ok {
				if sub, ok := claims["sub"].(string); ok {
					c.Set("userID", sub)
					c.Next()
					return
				}
			}
		}

		if rErr == nil {
			if newAccess, claims, ok := jwt.RefreshToken(refreshCookie.Value); ok {
				if sub, ok := claims["sub"].(string); ok {
					http.SetCookie(c.Writer, &http.Cookie{
						Name:     "access",
						Value:    newAccess,
						Path:     "/",
						HttpOnly: true,
						Secure:   true,
						SameSite: http.SameSiteLaxMode,
						Expires:  time.Now().Add(15 * time.Minute),
					})
					c.Set("userID", sub)
					c.Next()
					return
				}
			}
		}

		// редирект если не авторизован
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
	}
}