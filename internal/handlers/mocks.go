package handlers

import (
	"context"

	"github.com/anastasikalen/bookshelf.git/internal/models"
)

type mockBookRepo struct {
	CreateFn             func(ctx context.Context, book *models.Book) error
	GetByIDFn            func(ctx context.Context, id uint) (*models.Book, error)
	GetAllFn             func(ctx context.Context, limit, offset int) ([]models.Book, error)
	UpdateFn             func(ctx context.Context, book *models.Book) error
	DeleteFn             func(ctx context.Context, id uint) error
	MarkOutOfStockFn     func(ctx context.Context, id uint) error
	GetRecommendationsFn func(ctx context.Context) ([]models.Book, error)
}

func (m *mockBookRepo) Create(ctx context.Context, book *models.Book) error {
	return m.CreateFn(ctx, book)
}
func (m *mockBookRepo) GetByID(ctx context.Context, id uint) (*models.Book, error) {
	return m.GetByIDFn(ctx, id)
}
func (m *mockBookRepo) GetAll(ctx context.Context, limit, offset int) ([]models.Book, error) {
	return m.GetAllFn(ctx, limit, offset)
}
func (m *mockBookRepo) Update(ctx context.Context, book *models.Book) error {
	return m.UpdateFn(ctx, book)
}
func (m *mockBookRepo) Delete(ctx context.Context, id uint) error {
	return m.DeleteFn(ctx, id)
}
func (m *mockBookRepo) MarkOutOfStock(ctx context.Context, id uint) error {
	return m.MarkOutOfStockFn(ctx, id)
}
func (m *mockBookRepo) GetRecommendations(ctx context.Context) ([]models.Book, error) {
	return m.GetRecommendations(ctx)
}
