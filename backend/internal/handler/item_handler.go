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
//
//	@Summary		List items
//	@Description	Returns a paginated list of active items (soft-deleted excluded).
//	@Tags			items
//	@Produce		json
//	@Param			page		query		int	false	"Page number"	default(1)
//	@Param			page_size	query		int	false	"Page size (10, 20, 50)"	default(10)
//	@Success		200			{object}	domain.PaginatedItems
//	@Failure		500			{object}	map[string]string
//	@Router			/items [get]
func (h *ItemHandler) List(c fiber.Ctx) error {
	page, pageSize := parsePagination(c)
	result, err := h.uc.List(c.Context(), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

// Get returns a single item by UUID including its ingredient UUIDs.
//
//	@Summary		Get item
//	@Description	Returns a single item by UUID with its associated ingredient UUIDs.
//	@Tags			items
//	@Produce		json
//	@Param			uuid	path		string	true	"Item UUID"
//	@Success		200		{object}	domain.Item
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/items/{uuid} [get]
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

// Create inserts a new item with ingredient relationships inside a transaction.
//
//	@Summary		Create item
//	@Description	Creates a new item. Name must be unique. At least one ingredient is required. Ingredient UUIDs must exist.
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			body	body		object{name=string,price=number,status=int,ingredients=[]string}	true	"Item payload"
//	@Success		201		{object}	domain.Item
//	@Failure		400		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Router			/items [post]
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

// Update modifies an item and its ingredient relationships inside a transaction.
//
//	@Summary		Update item
//	@Description	Updates an item. Name must be unique excluding current record. Replaces all ingredient relationships.
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string																true	"Item UUID"
//	@Param			body	body		object{name=string,price=number,status=int,ingredients=[]string}	true	"Item payload"
//	@Success		200		{object}	domain.Item
//	@Failure		400		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Router			/items/{uuid} [put]
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

// Delete soft-deletes an item.
//
//	@Summary		Delete item
//	@Description	Soft-deletes an item (sets deleted_at). Note: tm_item_ingredient relationships are NOT deleted.
//	@Tags			items
//	@Param			uuid	path	string	true	"Item UUID"
//	@Success		204
//	@Failure		500	{object}	map[string]string
//	@Router			/items/{uuid} [delete]
func (h *ItemHandler) Delete(c fiber.Ctx) error {
	if err := h.uc.Delete(c.Context(), c.Params("uuid")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

