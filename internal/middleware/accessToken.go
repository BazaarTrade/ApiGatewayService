package middleware

import (
	"net/http"
	"strings"

	"github.com/BazaarTrade/ApiGatewayService/internal/tokenManager"
	"github.com/labstack/echo/v4"
)

func AccessTokenMiddleware(next echo.HandlerFunc, ACCESS_TOKEN_SECRET_PHRASE string) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "authorization header missing",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid authorization header format",
			})
		}

		accessToken := parts[1]

		isValid, err := tokenManager.ValidateAccessToken(accessToken, ACCESS_TOKEN_SECRET_PHRASE)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid access token",
			})
		}

		if !isValid {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "access token expired",
			})
		}

		return next(c)
	}
}
