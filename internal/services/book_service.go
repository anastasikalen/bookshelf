package services

import (
	"context"

	"github.com/anastasikalen/bookshelf.git/internal/models"
	"github.com/anastasikalen/bookshelf.git/internal/repository"
)

type BookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) MarkOutOfStock(ctx context.Context, id uint) error {
	return s.repo.MarkOutOfStock(ctx, id)
}

func (s *BookService) GetRecommendations(ctx context.Context) ([]models.Book, error) {
	return s.repo.GetRecommendations(ctx)
}
