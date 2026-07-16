package handler

import (
	"tablelink-backend/internal/domain"
	"tablelink-backend/internal/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type IngredientHandler struct {
	uc usecase.IngredientUsecase
}

func NewIngredientHandler(uc usecase.IngredientUsecase) *IngredientHandler {
	return &IngredientHandler{uc: uc}
}

func (h *IngredientHandler) Register(r fiber.Router) {
	grp := r.Group("/ingredients")
	grp.Get("/", h.List)
	grp.Get("/:uuid", h.Get)
	grp.Post("/", h.Create)
	grp.Put("/:uuid", h.Update)
	grp.Delete("/:uuid", h.Delete)
}

// List returns a paginated list of active ingredients.
// Query params: page (default 1), page_size (10, 20, or 50; default 10).
func (h *IngredientHandler) List(c fiber.Ctx) error {
	page, pageSize := parsePagination(c)
	result, err := h.uc.List(c.Context(), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *IngredientHandler) Get(c fiber.Ctx) error {
	ingredient, err := h.uc.Get(c.Context(), c.Params("uuid"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if ingredient == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "ingredient not found"})
	}
	return c.JSON(ingredient)
}

func (h *IngredientHandler) Create(c fiber.Ctx) error {
	var input struct {
		Name        string                 `json:"name"`
		CauseAlergy bool                   `json:"cause_alergy"`
		Type        domain.IngredientType  `json:"type"`
		Status      domain.RecordStatus    `json:"status"`
	}
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if input.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}
	ingredient := &domain.Ingredient{
		UUID: uuid.New().String(), Name: input.Name,
		CauseAlergy: input.CauseAlergy, Type: input.Type, Status: input.Status,
	}
	if err := h.uc.Create(c.Context(), ingredient); err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(ingredient)
}

func (h *IngredientHandler) Update(c fiber.Ctx) error {
	var input struct {
		Name        string                 `json:"name"`
		CauseAlergy bool                   `json:"cause_alergy"`
		Type        domain.IngredientType  `json:"type"`
		Status      domain.RecordStatus    `json:"status"`
	}
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if input.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}
	ingredient := &domain.Ingredient{
		UUID: c.Params("uuid"), Name: input.Name,
		CauseAlergy: input.CauseAlergy, Type: input.Type, Status: input.Status,
	}
	if err := h.uc.Update(c.Context(), ingredient); err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ingredient)
}

func (h *IngredientHandler) Delete(c fiber.Ctx) error {
	if err := h.uc.Delete(c.Context(), c.Params("uuid")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
