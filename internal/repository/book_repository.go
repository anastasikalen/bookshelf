package repository

import (
	"context"
	"errors"

	"github.com/anastasikalen/bookshelf.git/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepository interface {
	Create(ctx context.Context, book *models.Book) error
	GetByID(ctx context.Context, id uint) (*models.Book, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.Book, error)
	Update(ctx context.Context, book *models.Book) error
	Delete(ctx context.Context, id uint) error
	MarkOutOfStock(ctx context.Context, id uint) error
	GetRecommendations(ctx context.Context) ([]models.Book, error)
}

type bookRepository struct {
	db *pgxpool.Pool
}

func (r bookRepository) Create(ctx context.Context, book *models.Book) error {
	query := `
INSERT INTO books (title, author, year, isbn, out_of_stock, rating)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(
		ctx,
		query,
		book.Title,
		book.Author,
		book.Year,
		book.ISBN,
		book.OutOfStock,
		book.Rating,
	).Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)
}

func (r bookRepository) GetByID(ctx context.Context, id uint) (*models.Book, error) {
	query := `
SELECT id, title, author, year, isbn, out_of_stock, rating, created_at, updated_at
FROM books
WHERE id = $1`

	var book models.Book

	err := r.db.QueryRow(ctx, query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Year,
		&book.ISBN,
		&book.OutOfStock,
		&book.Rating,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r bookRepository) GetAll(ctx context.Context, limit, offset int) ([]models.Book, error) {
	query := `
SELECT id, title, author, year, isbn, out_of_stock, rating, created_at, updated_at
FROM books
ORDER BY id
LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Year,
			&book.ISBN,
			&book.OutOfStock,
			&book.Rating,
			&book.CreatedAt,
			&book.UpdatedAt,
		); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r bookRepository) Update(ctx context.Context, book *models.Book) error {
	query := `
UPDATE books
SET title = $1,
author = $2,
year = $3,
isbn = $4,
out_of_stock = $5,
rating = $6,
updated_at = NOW()
WHERE id = $7`
	cmd, err := r.db.Exec(ctx,
		query,
		book.Title,
		book.Author,
		book.Year,
		book.ISBN,
		book.OutOfStock,
		book.Rating,
		book.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("book not found")
	}
	return nil
}

func (r bookRepository) Delete(ctx context.Context, id uint) error {
	query := `
DELETE FROM books
WHERE id = $1`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("book not found")
	}
	return nil
}

func (r bookRepository) MarkOutOfStock(ctx context.Context, id uint) error {
	query := `
UPDATE books SET out_of_stock = true, updated_at = NOW()
		WHERE id = $1`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("book not found")
	}
	return nil
}

func (r bookRepository) GetRecommendations(ctx context.Context) ([]models.Book, error) {
	query := `
SELECT id, title, author, year, isbn, out_of_stock, rating, created_at, updated_at
FROM books
ORDER BY rating DESC NULLS LAST, year DESC
LIMIT 5`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var books []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Author,
			&b.Year,
			&b.ISBN,
			&b.OutOfStock,
			&b.Rating,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func NewBookRepository(db *pgxpool.Pool) BookRepository {
	return &bookRepository{db: db}
}
