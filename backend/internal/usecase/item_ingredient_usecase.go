package usecase

import (
	"context"

	"tablelink-backend/internal/domain"
	"tablelink-backend/internal/repository"
)

// ItemIngredientUsecase exposes read-only access to item-ingredient
// relationships. Create / Update / Delete are handled inside ItemUsecase
// within PostgreSQL transactions for consistency.
type ItemIngredientUsecase interface {
	ListByItem(ctx context.Context, itemUUID string) ([]domain.ItemIngredient, error)
}

type itemIngredientUC struct {
	repo repository.ItemIngredientRepository
}

func NewItemIngredientUsecase(repo repository.ItemIngredientRepository) ItemIngredientUsecase {
	return &itemIngredientUC{repo: repo}
}

func (uc *itemIngredientUC) ListByItem(ctx context.Context, itemUUID string) ([]domain.ItemIngredient, error) {
	return uc.repo.FindByItemUUID(ctx, itemUUID)
}

