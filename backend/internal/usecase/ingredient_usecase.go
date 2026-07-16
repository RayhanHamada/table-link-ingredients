package usecase

import (
	"context"
	"errors"
	"fmt"

	"tablelink-backend/internal/domain"
	"tablelink-backend/internal/repository"

	"github.com/jackc/pgx/v5"
)

// ---------------------------------------------------------------------------
// Interface
// ---------------------------------------------------------------------------

// IngredientUsecase defines business-logic operations for ingredients.
type IngredientUsecase interface {
	List(ctx context.Context, page, pageSize int) (*domain.PaginatedIngredients, error)
	Get(ctx context.Context, uuid string) (*domain.Ingredient, error)
	Create(ctx context.Context, i *domain.Ingredient) error
	Update(ctx context.Context, i *domain.Ingredient) error
	Delete(ctx context.Context, uuid string) error
}

// ---------------------------------------------------------------------------
// Implementation
// ---------------------------------------------------------------------------

type ingredientUC struct {
	repo repository.IngredientRepository
}

// NewIngredientUsecase wires the ingredient usecase.
func NewIngredientUsecase(repo repository.IngredientRepository) IngredientUsecase {
	return &ingredientUC{repo: repo}
}

func (uc *ingredientUC) List(ctx context.Context, page, pageSize int) (*domain.PaginatedIngredients, error) {
	data, total, err := uc.repo.FindAllPaginated(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &domain.PaginatedIngredients{
		Data:     data,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func (uc *ingredientUC) Get(ctx context.Context, uuid string) (*domain.Ingredient, error) {
	return uc.repo.FindByUUID(ctx, uuid)
}

// Create validates then inserts a new ingredient.
// Validation: name must be unique (ignoring soft-deleted records).
func (uc *ingredientUC) Create(ctx context.Context, i *domain.Ingredient) error {
	existing, err := uc.repo.FindByName(ctx, i.Name)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if existing != nil {
		return fmt.Errorf("ingredient name %q already exists", i.Name)
	}
	return uc.repo.Create(ctx, i)
}

// Update validates then modifies an ingredient.
// Validation: name must be unique, ignoring the current record and soft-deleted.
func (uc *ingredientUC) Update(ctx context.Context, i *domain.Ingredient) error {
	existing, err := uc.repo.FindByName(ctx, i.Name)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if existing != nil && existing.UUID != i.UUID {
		return fmt.Errorf("ingredient name %q already exists", i.Name)
	}
	return uc.repo.Update(ctx, i)
}

func (uc *ingredientUC) Delete(ctx context.Context, uuid string) error {
	return uc.repo.Delete(ctx, uuid)
}

