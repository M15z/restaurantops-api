package services

import (
	"context"
	"errors"

	"github.com/m15z/restaurantops-api/internal/models"
	"github.com/m15z/restaurantops-api/internal/repository"
)

var ErrInvalidInput = errors.New("invalid input")

type MenuService struct {
	categories *repository.CategoryRepository
	items      *repository.MenuItemRepository
}

func NewMenuService(categories *repository.CategoryRepository, items *repository.MenuItemRepository) *MenuService {
	return &MenuService{categories: categories, items: items}
}

func (s *MenuService) CreateCategory(ctx context.Context, name string) (int64, error) {
	if name == "" {
		return 0, ErrInvalidInput
	}

	return s.categories.Create(ctx, name)
}

func (s *MenuService) GetCategories(ctx context.Context) ([]models.Category, error) {
	return s.categories.GetAll(ctx)
}

// --- Menu items ---

func (s *MenuService) CreateItem(ctx context.Context, m models.MenuItem) (int64, error) {
	if err := validateItem(m); err != nil {
		return 0, err
	}
	return s.items.Create(ctx, m)
}

func (s *MenuService) GetItems(ctx context.Context) ([]models.MenuItem, error) {
	return s.items.GetAll(ctx)
}

func (s *MenuService) GetItem(ctx context.Context, id int64) (*models.MenuItem, error) {
	return s.items.GetByID(ctx, id)
}

func (s *MenuService) UpdateItem(ctx context.Context, m models.MenuItem) error {
	if err := validateItem(m); err != nil {
		return err
	}

	return s.items.Update(ctx, m)
}

func (s *MenuService) SetItemAvailability(ctx context.Context, id int64, available bool) error {
	return s.items.SetAvailable(ctx, id, available)
}

func validateItem(m models.MenuItem) error {
	if m.Name == "" || m.Price <= 0 || m.Category.ID <= 0 {
		return ErrInvalidInput
	}
	return nil
}
