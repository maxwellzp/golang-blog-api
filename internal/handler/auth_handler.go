package handler

import (
	"github.com/labstack/echo/v4"
	"log"
	"maxwellzp/blog-api/internal/service"
	"net/http"
)

type AuthHandler struct {
	AuthService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
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
		log.Printf("Error binding request: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	user, err := h.AuthService.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		log.Printf("Error registering user: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "something went wrong"})
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	var req loginRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("Error binding request: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	user, err := h.AuthService.Login(ctx, req.Email, req.Password)
	if err != nil {
		log.Printf("Error logining user: %v", err)
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "something went wrong"})
	}

	return c.JSON(http.StatusOK, user)
}
