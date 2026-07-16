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
//
//	@Summary		List ingredients
//	@Description	Returns a paginated list of active ingredients (soft-deleted excluded).
//	@Tags			ingredients
//	@Produce		json
//	@Param			page		query		int	false	"Page number"	default(1)
//	@Param			page_size	query		int	false	"Page size (10, 20, 50)"	default(10)
//	@Success		200			{object}	domain.PaginatedIngredients
//	@Failure		500			{object}	map[string]string
//	@Router			/ingredients [get]
func (h *IngredientHandler) List(c fiber.Ctx) error {
	page, pageSize := parsePagination(c)
	result, err := h.uc.List(c.Context(), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

// Get returns a single ingredient by UUID.
//
//	@Summary		Get ingredient
//	@Description	Returns a single ingredient by UUID.
//	@Tags			ingredients
//	@Produce		json
//	@Param			uuid	path		string	true	"Ingredient UUID"
//	@Success		200		{object}	domain.Ingredient
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/ingredients/{uuid} [get]
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

// Create inserts a new ingredient.
//
//	@Summary		Create ingredient
//	@Description	Creates a new ingredient. Name must be unique (excluding soft-deleted records).
//	@Tags			ingredients
//	@Accept			json
//	@Produce		json
//	@Param			body	body		object{name=string,cause_alergy=bool,type=int,status=int}	true	"Ingredient payload"
//	@Success		201		{object}	domain.Ingredient
//	@Failure		400		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Router			/ingredients [post]
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

// Update modifies an existing ingredient.
//
//	@Summary		Update ingredient
//	@Description	Updates an ingredient. Name must be unique excluding current record and soft-deleted.
//	@Tags			ingredients
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string														true	"Ingredient UUID"
//	@Param			body	body		object{name=string,cause_alergy=bool,type=int,status=int}	true	"Ingredient payload"
//	@Success		200		{object}	domain.Ingredient
//	@Failure		400		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Router			/ingredients/{uuid} [put]
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

// Delete soft-deletes an ingredient.
//
//	@Summary		Delete ingredient
//	@Description	Soft-deletes an ingredient (sets deleted_at).
//	@Tags			ingredients
//	@Param			uuid	path	string	true	"Ingredient UUID"
//	@Success		204
//	@Failure		500	{object}	map[string]string
//	@Router			/ingredients/{uuid} [delete]
func (h *IngredientHandler) Delete(c fiber.Ctx) error {
	if err := h.uc.Delete(c.Context(), c.Params("uuid")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
