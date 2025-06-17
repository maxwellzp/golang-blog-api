package middleware

import (
	"errors"
	"github.com/labstack/echo/v4"
)

var ErrUserIDNotFound = errors.New("user_id not found in context")

func GetUserID(c echo.Context) (int64, error) {
	val := c.Get(UserIDContextKey)
	if val == nil {
		return 0, ErrUserIDNotFound
	}

	userID, ok := val.(int64)
	if !ok {
		return 0, ErrUserIDNotFound
	}

	return userID, nil
}
