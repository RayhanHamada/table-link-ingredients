package handler

import (
	"tablelink-backend/internal/domain"
	"tablelink-backend/internal/usecase"

	"github.com/gofiber/fiber/v3"
)

type ItemHandler struct {
	uc usecase.ItemUsecase
}

func NewItemHandler(uc usecase.ItemUsecase) *ItemHandler {
	return &ItemHandler{uc: uc}
}

func (h *ItemHandler) Register(r fiber.Router) {
	grp := h.RegisterGroup(r)
	grp.Get("/", h.List)
	grp.Get("/:uuid", h.Get)
	grp.Post("/", h.Create)
	grp.Put("/:uuid", h.Update)
	grp.Delete("/:uuid", h.Delete)
}

func (h *ItemHandler) RegisterGroup(r fiber.Router) fiber.Router {
	return r.Group("/items")
}

// List returns a paginated list of active items.
// Query params: page (default 1), page_size (10, 20, or 50; default 10).
func (h *ItemHandler) List(c fiber.Ctx) error {
	page, pageSize := parsePagination(c)
	result, err := h.uc.List(c.Context(), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *ItemHandler) Get(c fiber.Ctx) error {
	item, err := h.uc.Get(c.Context(), c.Params("uuid"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if item == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "item not found"})
	}
	return c.JSON(item)
}

func (h *ItemHandler) Create(c fiber.Ctx) error {
	var input struct {
		Name        string             `json:"name"`
		Price       float64            `json:"price"`
		Status      domain.RecordStatus `json:"status"`
		Ingredients []string           `json:"ingredients"`
	}
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if input.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}
	if input.Price == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "price is required"})
	}
	if len(input.Ingredients) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "at least one ingredient is required"})
	}

	item, err := h.uc.Create(c.Context(), domain.ItemCreateInput{
		Name: input.Name, Price: input.Price, Status: input.Status, Ingredients: input.Ingredients,
	})
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *ItemHandler) Update(c fiber.Ctx) error {
	var input struct {
		Name        string              `json:"name"`
		Price       float64             `json:"price"`
		Status      domain.RecordStatus `json:"status"`
		Ingredients []string            `json:"ingredients"`
	}
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if input.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}
	if input.Price == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "price is required"})
	}
	if len(input.Ingredients) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "at least one ingredient is required"})
	}

	item, err := h.uc.Update(c.Context(), domain.ItemUpdateInput{
		UUID: c.Params("uuid"), Name: input.Name, Price: input.Price,
		Status: input.Status, Ingredients: input.Ingredients,
	})
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(item)
}

func (h *ItemHandler) Delete(c fiber.Ctx) error {
	if err := h.uc.Delete(c.Context(), c.Params("uuid")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

