package usecase

import (
	"context"

	"tablelink-backend/internal/domain"
	"tablelink-backend/internal/repository"
)

// ---------------------------------------------------------------------------
// Interface
// ---------------------------------------------------------------------------

// ItemIngredientUsecase defines business-logic operations for item-ingredient
// relationships.
type ItemIngredientUsecase interface {
	ListByItem(ctx context.Context, itemUUID string) ([]domain.ItemIngredient, error)
	ListByIngredient(ctx context.Context, ingredientUUID string) ([]domain.ItemIngredient, error)
	Create(ctx context.Context, rel *domain.ItemIngredient) error
	Delete(ctx context.Context, itemUUID, ingredientUUID string) error
}

// ---------------------------------------------------------------------------
// Implementation
// ---------------------------------------------------------------------------

type itemIngredientUC struct {
	repo repository.ItemIngredientRepository
}

// NewItemIngredientUsecase wires the item-ingredient usecase.
func NewItemIngredientUsecase(repo repository.ItemIngredientRepository) ItemIngredientUsecase {
	return &itemIngredientUC{repo: repo}
}

func (uc *itemIngredientUC) ListByItem(ctx context.Context, itemUUID string) ([]domain.ItemIngredient, error) {
	return uc.repo.FindByItemUUID(ctx, itemUUID)
}

func (uc *itemIngredientUC) ListByIngredient(ctx context.Context, ingredientUUID string) ([]domain.ItemIngredient, error) {
	return uc.repo.FindByIngredientUUID(ctx, ingredientUUID)
}

func (uc *itemIngredientUC) Create(ctx context.Context, rel *domain.ItemIngredient) error {
	return uc.repo.Create(ctx, rel)
}

func (uc *itemIngredientUC) Delete(ctx context.Context, itemUUID, ingredientUUID string) error {
	return uc.repo.Delete(ctx, itemUUID, ingredientUUID)
}
