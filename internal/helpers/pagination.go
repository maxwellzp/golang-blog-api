package helpers

import (
	"github.com/labstack/echo/v4"
	"strconv"
)

type Pagination struct {
	Page   int
	Limit  int
	Offset int
}

func GetPagination(c echo.Context) Pagination {
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	return Pagination{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}
