package handler

import (
	"github.com/labstack/echo/v4"
	"log"
	"maxwellzp/blog-api/internal/service"
	"net/http"
	"strconv"
)

type BlogHandler struct {
	BlogService service.BlogService
}

func NewBlogHandler(blogService service.BlogService) *BlogHandler {
	return &BlogHandler{BlogService: blogService}
}

type blogRequest struct {
	UserID  int64  `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *BlogHandler) Create(c echo.Context) error {
	var req blogRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("Error binding blog request: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	// c.Request().Context() extracts the context.Context from the incoming HTTP request.
	// This context includes: Timeout/cancel signals from the client.
	blog, err := h.BlogService.Create(c.Request().Context(), req.UserID, req.Title, req.Content)
	if err != nil {
		log.Printf("Error creating blog: %v", err)
		return c.JSON(http.StatusBadGateway, echo.Map{"error": "error creating blog"})
	}
	return c.JSON(http.StatusCreated, blog)
}

func (h *BlogHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Error parsing id param: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	blog, err := h.BlogService.GetByID(c.Request().Context(), id)
	if err != nil {
		log.Printf("Error getting blog: %v", err)
		return c.JSON(http.StatusNotFound, echo.Map{"error": "blog not found"})
	}
	return c.JSON(http.StatusOK, blog)
}

func (h *BlogHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Error parsing id param: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	var req blogRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("Error binding blog request: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	err = h.BlogService.Update(c.Request().Context(), id, req.Title, req.Content)
	if err != nil {
		log.Printf("Error updating blog: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "error updating blog"})
	}
	return c.NoContent(http.StatusOK)
}

func (h *BlogHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Error parsing id param: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	if err := h.BlogService.Delete(c.Request().Context(), id); err != nil {
		log.Printf("Error deleting blog: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error deleting blog"})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *BlogHandler) List(c echo.Context) error {
	blogs, err := h.BlogService.List(c.Request().Context())
	if err != nil {
		log.Printf("Error listing blogs: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error listing blogs"})
	}
	return c.JSON(http.StatusOK, blogs)
}
