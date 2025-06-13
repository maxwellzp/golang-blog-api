package handler

import (
	"github.com/labstack/echo/v4"
	"maxwellzp/blog-api/internal/service"
	"net/http"
	"strconv"
)

type CommentHandler struct {
	CommentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{CommentService: commentService}
}

type commentRequest struct {
	UserID  int64  `json:"user_id"`
	BlogID  int64  `json:"blog_id"`
	Content string `json:"content"`
}

func (h *CommentHandler) Create(c echo.Context) error {
	var req commentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	comment, err := h.CommentService.Create(c.Request().Context(), req.UserID, req.BlogID, req.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	comment, err := h.CommentService.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "comment not found"})
	}
	return c.JSON(http.StatusOK, comment)
}

func (h *CommentHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	var req commentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	if err := h.CommentService.Update(c.Request().Context(), id, req.Content); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

func (h *CommentHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	if err := h.CommentService.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

func (h *CommentHandler) ListByBlogID(c echo.Context) error {
	blogID, err := strconv.ParseInt(c.Param("blog_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid blog id"})
	}

	comments, err := h.CommentService.ListByBlogID(c.Request().Context(), blogID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, comments)
}
