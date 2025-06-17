package handler

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"maxwellzp/blog-api/internal/helpers"
	"maxwellzp/blog-api/internal/middleware"
	"maxwellzp/blog-api/internal/service"
	"net/http"
	"strconv"
)

type CommentHandler struct {
	CommentService service.CommentService
	Logger         *zap.SugaredLogger
}

func NewCommentHandler(commentService service.CommentService, logger *zap.SugaredLogger) *CommentHandler {
	return &CommentHandler{CommentService: commentService, Logger: logger}
}

type commentRequest struct {
	BlogID  int64  `json:"blog_id"`
	Content string `json:"content"`
}

func (h *CommentHandler) Create(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	var req commentRequest
	if err := c.Bind(&req); err != nil {
		h.Logger.Errorw("Error binding comment create request",
			"err", err,
			"user_id", userID,
			"blog_id", req.BlogID,
			"content", helpers.TruncateString(req.Content, 100),
			"status", http.StatusBadRequest,
		)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	comment, err := h.CommentService.Create(c.Request().Context(), userID, req.BlogID, req.Content)
	if err != nil {
		h.Logger.Errorw("Error creating comment",
			"err", err,
			"user_id", userID,
			"blog_id", req.BlogID,
			"content", helpers.TruncateString(req.Content, 100),
			"status", http.StatusBadRequest,
		)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	h.Logger.Infow("Comment created successfully",
		"comment_id", comment.ID,
		"status", http.StatusCreated,
	)
	return c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) GetByID(c echo.Context) error {
	rawID := c.Param("id")
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		h.Logger.Errorw("Error parsing comment id param in GetByID",
			"comment_id", rawID,
			"error", err,
			"status", http.StatusBadRequest,
		)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	comment, err := h.CommentService.GetByID(c.Request().Context(), id)
	if err != nil {
		h.Logger.Errorw("Failed to get comment by id",
			"comment_id", id,
			"error", err,
			"status", http.StatusNotFound,
		)
		return c.JSON(http.StatusNotFound, echo.Map{"error": "comment not found"})
	}

	h.Logger.Infow("Comment found successfully",
		"comment_id", comment.ID,
		"status", http.StatusOK,
	)
	return c.JSON(http.StatusOK, comment)
}

func (h *CommentHandler) Update(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	rawID := c.Param("id")
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		h.Logger.Errorw("Error parsing comment id param in Update",
			"comment_id", rawID,
			"error", err,
			"user_id", userID,
			"status", http.StatusBadRequest,
		)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	isOwner, err := h.CommentService.IsOwner(c.Request().Context(), id, userID)
	if err != nil {
		h.Logger.Errorw("Error checking comment ownership", "comment_id", id, "user_id", userID, "error", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}
	if !isOwner {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "you are not allowed to modify this comment"})
	}

	var req commentRequest
	if err := c.Bind(&req); err != nil {
		h.Logger.Errorw("Error binding comment update request",
			"comment_id", id,
			"error", err,
			"user_id", userID,
			"blog_id", req.BlogID,
			"content", helpers.TruncateString(req.Content, 100),
			"status", http.StatusBadRequest,
		)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	if err := h.CommentService.Update(c.Request().Context(), id, req.Content); err != nil {
		h.Logger.Errorw("Error updating comment",
			"comment_id", id,
			"error", err,
			"user_id", userID,
			"blog_id", req.BlogID,
			"content", helpers.TruncateString(req.Content, 100),
			"status", http.StatusBadRequest,
		)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	h.Logger.Infow("Comment updated successfully",
		"comment_id", id,
		"status", http.StatusOK,
	)
	return c.NoContent(http.StatusOK)
}

func (h *CommentHandler) Delete(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	rawID := c.Param("id")
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		h.Logger.Errorw("Error parsing comment id param in Delete",
			"comment_id", rawID,
			"error", err,
			"user_id", userID,
			"status", http.StatusBadRequest,
		)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	isOwner, err := h.CommentService.IsOwner(c.Request().Context(), id, userID)
	if err != nil {
		h.Logger.Errorw("Error checking comment ownership", "comment_id", id, "user_id", userID, "error", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}
	if !isOwner {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "you are not allowed to delete this comment"})
	}

	if err := h.CommentService.Delete(c.Request().Context(), id); err != nil {
		h.Logger.Errorw("Error deleting comment",
			"comment_id", id,
			"error", err,
			"user_id", userID,
			"status", http.StatusInternalServerError,
		)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	h.Logger.Infow("Comment deleted successfully",
		"comment_id", id,
		"status", http.StatusNoContent,
	)
	return c.NoContent(http.StatusNoContent)
}

func (h *CommentHandler) ListByBlogID(c echo.Context) error {
	rawBlogID := c.Param("blog_id")
	blogID, err := strconv.ParseInt(rawBlogID, 10, 64)
	if err != nil {
		h.Logger.Errorw("Error parsing blog_id param in ListByBlogID",
			"blog_id", rawBlogID,
			"error", err,
			"status", http.StatusBadRequest,
		)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid blog id"})
	}

	comments, err := h.CommentService.ListByBlogID(c.Request().Context(), blogID)
	if err != nil {
		h.Logger.Errorw("Error listing comments",
			"blog_id", blogID,
			"error", err,
			"status", http.StatusInternalServerError,
		)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	h.Logger.Infow("Comments listed successfully",
		"comment_count", len(comments),
		"status", http.StatusOK,
		"blog_id", blogID,
	)
	return c.JSON(http.StatusOK, comments)
}
