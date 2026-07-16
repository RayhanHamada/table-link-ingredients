package handler

import (
	"tablelink-backend/internal/usecase"

	"github.com/gofiber/fiber/v3"
)

// ItemIngredientHandler provides read-only access to item-ingredient
// relationships. All write operations are handled inside ItemUsecase
// via transactions.
type ItemIngredientHandler struct {
	uc usecase.ItemIngredientUsecase
}

func NewItemIngredientHandler(uc usecase.ItemIngredientUsecase) *ItemIngredientHandler {
	return &ItemIngredientHandler{uc: uc}
}

// Register mounts read-only nested routes under /items.
//   GET /items/:uuid/ingredients → list ingredient ids for an item
func (h *ItemIngredientHandler) Register(r fiber.Router) {
	r.Get("/:uuid/ingredients", h.ListByItem)
}

func (h *ItemIngredientHandler) ListByItem(c fiber.Ctx) error {
	rels, err := h.uc.ListByItem(c.Context(), c.Params("uuid"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(rels)
}

