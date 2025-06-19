package handler

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"maxwellzp/blog-api/internal/service"
	"maxwellzp/blog-api/internal/validation"
	"net/http"
)

type AuthHandler struct {
	AuthService service.AuthService
	Logger      *zap.SugaredLogger
	Validator   *validation.Validator
}

func NewAuthHandler(
	authService service.AuthService,
	logger *zap.SugaredLogger,
	validator *validation.Validator,
) *AuthHandler {
	return &AuthHandler{AuthService: authService, Logger: logger, Validator: validator}
}

type registerRequest struct {
	Username string `json:"username" validate:"required,min=5,max=30,alphanumunicode"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=12,max=40,containsuppercase,containslowercase,containsnumber,containsspecial"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=40"`
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
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid registration request"})
	}

	if fieldErrors := h.Validator.ValidateStruct(&req); fieldErrors != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "validation failed",
			"fields": fieldErrors,
		})
	}

	user, err := h.AuthService.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		h.Logger.Errorw("Error registering user",
			"error", err,
			"email", req.Email,
			"username", req.Username,
			"status", http.StatusBadRequest,
		)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid registration request"})
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

	if fieldErrors := h.Validator.ValidateStruct(&req); fieldErrors != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "validation failed",
			"fields": fieldErrors,
		})
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
