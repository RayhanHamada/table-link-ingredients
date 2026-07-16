package handler

import (
	"tablelink-backend/internal/domain"
	"tablelink-backend/internal/usecase"

	"github.com/gofiber/fiber/v3"
)

// ItemIngredientHandler exposes HTTP endpoints for tm_item_ingredient.
type ItemIngredientHandler struct {
	uc usecase.ItemIngredientUsecase
}

// NewItemIngredientHandler wires the handler with its usecase dependency.
func NewItemIngredientHandler(uc usecase.ItemIngredientUsecase) *ItemIngredientHandler {
	return &ItemIngredientHandler{uc: uc}
}

// ---------------------------------------------------------------------------
// Route registration
// ---------------------------------------------------------------------------

// Register mounts item-ingredient relationship routes on the provided router.
func (h *ItemIngredientHandler) Register(r fiber.Router) {
	// Nested under items for clarity:
	//   GET  /items/:uuid/ingredients   -> list ingredients for an item
	//   POST /items/:uuid/ingredients   -> link an ingredient to an item
	r.Get("/:uuid/ingredients", h.ListByItem)
	r.Post("/:uuid/ingredients", h.Create)
	r.Delete("/:uuid/ingredients/:ingredient_uuid", h.Delete)
}

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

// ListByItem returns all ingredients linked to an item.
func (h *ItemIngredientHandler) ListByItem(c fiber.Ctx) error {
	itemUUID := c.Params("uuid")
	rels, err := h.uc.ListByItem(c.Context(), itemUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(rels)
}

// Create links an ingredient to an item.
func (h *ItemIngredientHandler) Create(c fiber.Ctx) error {
	itemUUID := c.Params("uuid")

	var input struct {
		UUIDIngredient string `json:"uuid_ingredient"`
	}
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	rel := &domain.ItemIngredient{
		UUIDItem:       itemUUID,
		UUIDIngredient: input.UUIDIngredient,
	}

	if err := h.uc.Create(c.Context(), rel); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(rel)
}

// Delete removes an ingredient from an item.
func (h *ItemIngredientHandler) Delete(c fiber.Ctx) error {
	itemUUID := c.Params("uuid")
	ingredientUUID := c.Params("ingredient_uuid")

	if err := h.uc.Delete(c.Context(), itemUUID, ingredientUUID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
