package usecase

import (
	"context"

	"tablelink-backend/internal/domain"
	"tablelink-backend/internal/repository"
)

// ---------------------------------------------------------------------------
// Interface
// ---------------------------------------------------------------------------

// IngredientUsecase defines business-logic operations for ingredients.
type IngredientUsecase interface {
	List(ctx context.Context) ([]domain.Ingredient, error)
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

func (uc *ingredientUC) List(ctx context.Context) ([]domain.Ingredient, error) {
	return uc.repo.FindAll(ctx)
}

func (uc *ingredientUC) Get(ctx context.Context, uuid string) (*domain.Ingredient, error) {
	return uc.repo.FindByUUID(ctx, uuid)
}

func (uc *ingredientUC) Create(ctx context.Context, i *domain.Ingredient) error {
	return uc.repo.Create(ctx, i)
}

func (uc *ingredientUC) Update(ctx context.Context, i *domain.Ingredient) error {
	return uc.repo.Update(ctx, i)
}

func (uc *ingredientUC) Delete(ctx context.Context, uuid string) error {
	return uc.repo.Delete(ctx, uuid)
}
