package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

var validPageSizes = map[int]bool{10: true, 20: true, 50: true}

// parsePagination extracts page and page_size from Fiber query parameters.
// Defaults: page = 1, page_size = 10. Clamps page_size to 10, 20, or 50.
func parsePagination(c fiber.Ctx) (page int, pageSize int) {
	page = 1
	pageSize = 10

	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && validPageSizes[v] {
			pageSize = v
		}
	}
	return
}
