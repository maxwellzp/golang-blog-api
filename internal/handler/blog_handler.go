package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"maxwellzp/blog-api/internal/helpers"
	"maxwellzp/blog-api/internal/middleware"
	"maxwellzp/blog-api/internal/service"
	"maxwellzp/blog-api/internal/validation"
	"net/http"
	"strconv"
)

type BlogHandler struct {
	BlogService service.BlogService
	Logger      *zap.SugaredLogger
	Validator   *validation.Validator
}

func NewBlogHandler(
	blogService service.BlogService,
	logger *zap.SugaredLogger,
	validator *validation.Validator,
) *BlogHandler {
	return &BlogHandler{BlogService: blogService, Logger: logger, Validator: validator}
}

type blogRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=100"`
	Content string `json:"content" validate:"required,min=10"`
}

func (h *BlogHandler) Create(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	var req blogRequest
	if err := c.Bind(&req); err != nil {
		h.Logger.Errorw("Error binding blog create request",
			"error", err,
			"user_id", userID,
			"title", helpers.TruncateString(req.Title, 100),
			"content", helpers.TruncateString(req.Content, 100),
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

	// c.Request().Context() extracts the context.Context from the incoming HTTP request.
	// This context includes: Timeout/cancel signals from the client.
	blog, err := h.BlogService.Create(c.Request().Context(), userID, req.Title, req.Content)
	if err != nil {
		h.Logger.Errorw("Error creating blog",
			"error", err,
			"user_id", userID,
			"title", helpers.TruncateString(req.Title, 100),
			"content", helpers.TruncateString(req.Content, 100),
			"status", http.StatusBadRequest,
		)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "error creating blog"})
	}

	h.Logger.Infow("Blog created successfully",
		"blog_id", blog.ID,
		"status", http.StatusCreated,
	)
	return c.JSON(http.StatusCreated, blog)
}

func (h *BlogHandler) GetByID(c echo.Context) error {
	rawID := c.Param("id")
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		h.Logger.Errorw("Error parsing id param in GetByID",
			"blog_id", rawID,
			"error", err,
			"status", http.StatusBadRequest,
		)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	blog, err := h.BlogService.GetByID(c.Request().Context(), id)
	if err != nil {
		h.Logger.Errorw("Failed to get blog by id",
			"blog_id", id,
			"error", err,
			"status", http.StatusNotFound,
		)
		return c.JSON(http.StatusNotFound, echo.Map{"error": "blog not found"})
	}

	h.Logger.Infow("Blog found successfully",
		"blog_id", blog.ID,
		"status", http.StatusOK,
	)
	return c.JSON(http.StatusOK, blog)
}

func (h *BlogHandler) Update(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	rawID := c.Param("id")
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		h.Logger.Errorw("Error parsing id param in Update",
			"blog_id", rawID,
			"error", err,
			"user_id", userID,
			"status", http.StatusBadRequest,
		)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	isOwner, err := h.BlogService.IsOwner(c.Request().Context(), id, userID)
	if err != nil {
		if errors.Is(err, service.ErrBlogNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "blog not found"})
		}
		h.Logger.Errorw("Error checking blog ownership", "blog_id", id, "user_id", userID, "error", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}
	if !isOwner {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "you are not allowed to modify this blog"})
	}

	var req blogRequest
	if err := c.Bind(&req); err != nil {
		h.Logger.Errorw("Error binding blog update request",
			"blog_id", id,
			"error", err,
			"user_id", userID,
			"title", helpers.TruncateString(req.Title, 100),
			"content", helpers.TruncateString(req.Content, 100),
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

	err = h.BlogService.Update(c.Request().Context(), id, req.Title, req.Content)
	if err != nil {
		h.Logger.Errorw("Error updating blog",
			"blog_id", id,
			"error", err,
			"user_id", userID,
			"title", helpers.TruncateString(req.Title, 100),
			"content", helpers.TruncateString(req.Content, 100),
			"status", http.StatusBadRequest,
		)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "error updating blog"})
	}

	h.Logger.Infow("Blog updated successfully",
		"blog_id", id,
		"status", http.StatusOK,
	)
	return c.NoContent(http.StatusOK)
}

func (h *BlogHandler) Delete(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	rawID := c.Param("id")
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		h.Logger.Errorw("Error parsing id param in Delete",
			"blog_id", rawID,
			"error", err,
			"user_id", userID,
			"status", http.StatusBadRequest,
		)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	isOwner, err := h.BlogService.IsOwner(c.Request().Context(), id, userID)
	if err != nil {
		if errors.Is(err, service.ErrBlogNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "blog not found"})
		}
		h.Logger.Errorw("Error checking blog ownership", "blog_id", id, "user_id", userID, "error", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}
	if !isOwner {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "you are not allowed to delete this blog"})
	}

	if err := h.BlogService.Delete(c.Request().Context(), id); err != nil {
		h.Logger.Errorw("Error deleting blog",
			"blog_id", id,
			"error", err,
			"user_id", userID,
			"status", http.StatusInternalServerError,
		)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error deleting blog"})
	}

	h.Logger.Infow("Blog deleted successfully",
		"blog_id", id,
		"status", http.StatusNoContent,
	)
	return c.NoContent(http.StatusNoContent)
}

func (h *BlogHandler) List(c echo.Context) error {
	pagination := helpers.GetPagination(c)
	blogs, err := h.BlogService.List(c.Request().Context(), pagination.Limit, pagination.Offset)
	if err != nil {
		h.Logger.Errorw("Error listing blogs",
			"error", err,
			"status", http.StatusInternalServerError,
		)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error listing blogs"})
	}

	h.Logger.Infow("Blogs listed successfully",
		"blog_count", len(blogs),
		"status", http.StatusOK,
	)
	return c.JSON(http.StatusOK, blogs)
}
