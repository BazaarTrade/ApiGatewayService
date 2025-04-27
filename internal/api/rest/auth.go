package rest

import (
	"net/http"
	"strconv"
	"time"

	errorHandler "github.com/BazaarTrade/ApiGatewayService/internal/api/gRPC/errors"
	"github.com/BazaarTrade/ApiGatewayService/internal/models"
	"github.com/BazaarTrade/ApiGatewayService/internal/tokenManager"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/status"
)

func (s *Server) register(c echo.Context) error {
	var req models.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request format",
		})
	}

	if err := s.aClient.Register(req.Email, req.Password); err != nil {
		if st, ok := status.FromError(err); ok {
			return errorHandler.HandleGRPCError(c, st)
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "failed to register user",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "user registered successfully",
	})
}

func (s *Server) login(c echo.Context) error {
	var (
		req          models.LoginRequest
		userID       int
		accessToken  string
		refreshToken string
		err          error
	)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request format",
		})
	}

	userID, refreshToken, err = s.aClient.Login(req.Email, req.Password)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return errorHandler.HandleGRPCError(c, st)
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "failed to login user",
		})
	}

	accessToken, err = tokenManager.GenerateAccessToken(userID, s.accessTokenSecretPhrase)
	if err != nil {
		s.logger.Error("failed to generate access token", "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate access token",
		})
	}

	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(30 * 24 * time.Hour), // 30 days expiration
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteStrictMode,
	})

	c.Response().Header().Set("Authorization", "Bearer "+accessToken)

	return c.JSON(http.StatusOK, map[string]int{
		"userID": userID,
	})
}

func (s *Server) logout(c echo.Context) error {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "failed to get refresh token from cookie",
		})
	}

	if refreshToken == nil || refreshToken.Value == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "fialed to find refresh token",
		})
	}

	if err := s.aClient.Logout(refreshToken.Value); err != nil {
		if st, ok := status.FromError(err); ok {
			return errorHandler.HandleGRPCError(c, st)
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "failed to logout user",
		})
	}

	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteStrictMode,
	})

	return c.JSON(http.StatusOK, map[string]string{
		"message": "user logged out successfully",
	})
}

func (s *Server) refreshAccessToken(c echo.Context) error {
	userIDString := c.Param("userID")
	if userIDString == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "userID is required",
		})
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid userID format",
		})
	}

	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "failed to get refresh token from cookie",
		})
	}

	if refreshToken == nil || refreshToken.Value == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "fialed to find refresh token",
		})
	}

	isValid, newRefreshToken, err := s.aClient.IsRefreshTokenValid(userID, refreshToken.Value)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return errorHandler.HandleGRPCError(c, st)
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "failed to refresh token",
		})
	}

	if !isValid {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid refreshToken",
		})
	}

	accessToken, err := tokenManager.GenerateAccessToken(userID, s.accessTokenSecretPhrase)
	if err != nil {
		s.logger.Error("failed to generate access token", "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate access token",
		})
	}

	if newRefreshToken != "" {
		c.SetCookie(&http.Cookie{
			Name:     "refreshToken",
			Value:    newRefreshToken,
			Path:     "/",
			Expires:  time.Now().Add(30 * 24 * time.Hour), // 30 days expiration
			HttpOnly: true,
			Secure:   false, // Set to true in production
			SameSite: http.SameSiteStrictMode,
		})
	}

	c.Response().Header().Set("Authorization", "Bearer "+accessToken)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "token refreshed successfully",
	})
}

func (s *Server) changePassword(c echo.Context) error {
	var req models.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request format",
		})
	}

	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "failed to get refresh token from cookie",
		})
	}

	if refreshToken == nil || refreshToken.Value == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "fialed to find refresh token",
		})
	}

	if err := s.aClient.ChangePassword(req.Email, req.OldPassword, req.NewPassword, refreshToken.Value); err != nil {
		if st, ok := status.FromError(err); ok {
			return errorHandler.HandleGRPCError(c, st)
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "failed to change password",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "password changed successfully",
	})
}
