package handler

import (
	"tablelink-backend/internal/domain"
	"tablelink-backend/internal/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// ItemHandler exposes HTTP endpoints for tm_item.
type ItemHandler struct {
	uc usecase.ItemUsecase
}

// NewItemHandler wires the handler with its usecase dependency.
func NewItemHandler(uc usecase.ItemUsecase) *ItemHandler {
	return &ItemHandler{uc: uc}
}

// ---------------------------------------------------------------------------
// Route registration
// ---------------------------------------------------------------------------

// Register mounts item routes on the provided Fiber router.
func (h *ItemHandler) Register(r fiber.Router) {
	grp := h.RegisterGroup(r)
	grp.Get("/", h.List)
	grp.Get("/:uuid", h.Get)
	grp.Post("/", h.Create)
	grp.Put("/:uuid", h.Update)
	grp.Delete("/:uuid", h.Delete)
}

// RegisterGroup creates the /items group so that nested handlers (e.g.
// item-ingredient) can be mounted underneath it.
func (h *ItemHandler) RegisterGroup(r fiber.Router) fiber.Router {
	return r.Group("/items")
}

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

// List returns all active items.
func (h *ItemHandler) List(c fiber.Ctx) error {
	items, err := h.uc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(items)
}

// Get returns a single item by UUID.
func (h *ItemHandler) Get(c fiber.Ctx) error {
	id := c.Params("uuid")
	item, err := h.uc.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if item == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "item not found",
		})
	}
	return c.JSON(item)
}

// Create inserts a new item.
func (h *ItemHandler) Create(c fiber.Ctx) error {
	var input struct {
		Name   string             `json:"name"`
		Price  float64            `json:"price"`
		Status domain.RecordStatus `json:"status"`
	}
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	item := &domain.Item{
		UUID:   uuid.New().String(),
		Name:   input.Name,
		Price:  input.Price,
		Status: input.Status,
	}

	if err := h.uc.Create(c.Context(), item); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

// Update modifies an existing item.
func (h *ItemHandler) Update(c fiber.Ctx) error {
	id := c.Params("uuid")

	var input struct {
		Name   string              `json:"name"`
		Price  float64             `json:"price"`
		Status domain.RecordStatus `json:"status"`
	}
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	item := &domain.Item{
		UUID:   id,
		Name:   input.Name,
		Price:  input.Price,
		Status: input.Status,
	}

	if err := h.uc.Update(c.Context(), item); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(item)
}

// Delete soft-deletes an item.
func (h *ItemHandler) Delete(c fiber.Ctx) error {
	id := c.Params("uuid")
	if err := h.uc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
