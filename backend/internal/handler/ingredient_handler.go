package handler

import (
	"tablelink-backend/internal/domain"
	"tablelink-backend/internal/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// IngredientHandler exposes HTTP endpoints for tm_ingredient.
type IngredientHandler struct {
	uc usecase.IngredientUsecase
}

// NewIngredientHandler wires the handler with its usecase dependency.
func NewIngredientHandler(uc usecase.IngredientUsecase) *IngredientHandler {
	return &IngredientHandler{uc: uc}
}

// ---------------------------------------------------------------------------
// Route registration
// ---------------------------------------------------------------------------

// Register mounts ingredient routes on the provided Fiber router.
func (h *IngredientHandler) Register(r fiber.Router) {
	grp := r.Group("/ingredients")
	grp.Get("/", h.List)
	grp.Get("/:uuid", h.Get)
	grp.Post("/", h.Create)
	grp.Put("/:uuid", h.Update)
	grp.Delete("/:uuid", h.Delete)
}

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

// List returns all active ingredients.
func (h *IngredientHandler) List(c fiber.Ctx) error {
	ingredients, err := h.uc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(ingredients)
}

// Get returns a single ingredient by UUID.
func (h *IngredientHandler) Get(c fiber.Ctx) error {
	id := c.Params("uuid")
	ingredient, err := h.uc.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if ingredient == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "ingredient not found",
		})
	}
	return c.JSON(ingredient)
}

// Create inserts a new ingredient.
func (h *IngredientHandler) Create(c fiber.Ctx) error {
	var input struct {
		Name        string              `json:"name"`
		CauseAlergy bool                `json:"cause_alergy"`
		Type        domain.IngredientType `json:"type"`
		Status      domain.RecordStatus   `json:"status"`
	}
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	ingredient := &domain.Ingredient{
		UUID:        uuid.New().String(),
		Name:        input.Name,
		CauseAlergy: input.CauseAlergy,
		Type:        input.Type,
		Status:      input.Status,
	}

	if err := h.uc.Create(c.Context(), ingredient); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(ingredient)
}

// Update modifies an existing ingredient.
func (h *IngredientHandler) Update(c fiber.Ctx) error {
	id := c.Params("uuid")

	var input struct {
		Name        string                `json:"name"`
		CauseAlergy bool                  `json:"cause_alergy"`
		Type        domain.IngredientType `json:"type"`
		Status      domain.RecordStatus   `json:"status"`
	}
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	ingredient := &domain.Ingredient{
		UUID:        id,
		Name:        input.Name,
		CauseAlergy: input.CauseAlergy,
		Type:        input.Type,
		Status:      input.Status,
	}

	if err := h.uc.Update(c.Context(), ingredient); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(ingredient)
}

// Delete soft-deletes an ingredient.
func (h *IngredientHandler) Delete(c fiber.Ctx) error {
	id := c.Params("uuid")
	if err := h.uc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
