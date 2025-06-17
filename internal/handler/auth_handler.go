package handler

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"maxwellzp/blog-api/internal/service"
	"net/http"
)

type AuthHandler struct {
	AuthService service.AuthService
	Logger      *zap.SugaredLogger
}

func NewAuthHandler(authService service.AuthService, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{AuthService: authService, Logger: logger}
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c echo.Context) error {
	ctx := c.Request().Context()

	var req registerRequest
	if err := c.Bind(&req); err != nil {
		h.Logger.Errorw("Error binding register request",
			"error", err,
			"email", req.Email,
			"username", req.Username,
			"status", http.StatusBadRequest,
		)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	user, err := h.AuthService.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		h.Logger.Errorw("Error registering user",
			"error", err,
			"email", req.Email,
			"username", req.Username,
			"status", http.StatusBadRequest,
		)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "something went wrong"})
	}

	h.Logger.Infow("Successfully registered user",
		"user", user,
		"status", http.StatusCreated,
	)
	return c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	var req loginRequest
	if err := c.Bind(&req); err != nil {
		h.Logger.Errorw(
			"Error binding login request",
			"error", err,
			"email", req.Email,
			"status", http.StatusBadRequest,
		)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	user, token, err := h.AuthService.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.Logger.Errorw("Error logging user",
			"error", err,
			"email", req.Email,
			"status", http.StatusUnauthorized,
		)
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "something went wrong"})
	}

	h.Logger.Infow("User logged in",
		"user", user,
		"status", http.StatusOK,
	)
	return c.JSON(http.StatusOK, echo.Map{
		"user":  user,
		"token": token,
	})
}
