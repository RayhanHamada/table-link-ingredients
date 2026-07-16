package usecase

import (
	"context"

	"tablelink-backend/internal/domain"
	"tablelink-backend/internal/repository"
)

// ---------------------------------------------------------------------------
// Interface
// ---------------------------------------------------------------------------

// ItemUsecase defines business-logic operations for items.
type ItemUsecase interface {
	List(ctx context.Context) ([]domain.Item, error)
	Get(ctx context.Context, uuid string) (*domain.Item, error)
	Create(ctx context.Context, item *domain.Item) error
	Update(ctx context.Context, item *domain.Item) error
	Delete(ctx context.Context, uuid string) error
}

// ---------------------------------------------------------------------------
// Implementation
// ---------------------------------------------------------------------------

type itemUC struct {
	repo repository.ItemRepository
}

// NewItemUsecase wires the item usecase.
func NewItemUsecase(repo repository.ItemRepository) ItemUsecase {
	return &itemUC{repo: repo}
}

func (uc *itemUC) List(ctx context.Context) ([]domain.Item, error) {
	return uc.repo.FindAll(ctx)
}

func (uc *itemUC) Get(ctx context.Context, uuid string) (*domain.Item, error) {
	return uc.repo.FindByUUID(ctx, uuid)
}

func (uc *itemUC) Create(ctx context.Context, i *domain.Item) error {
	return uc.repo.Create(ctx, i)
}

func (uc *itemUC) Update(ctx context.Context, i *domain.Item) error {
	return uc.repo.Update(ctx, i)
}

func (uc *itemUC) Delete(ctx context.Context, uuid string) error {
	return uc.repo.Delete(ctx, uuid)
}
