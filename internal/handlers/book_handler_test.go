package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anastasikalen/bookshelf.git/internal/models"
	"github.com/anastasikalen/bookshelf.git/internal/services"
)

func TestCreateBook(t *testing.T) {
	repo := &mockBookRepo{
		CreateFn: func(ctx context.Context, book *models.Book) error {
			book.ID = 1
			return nil
		},
	}
	service := &services.BookService{}

	handler := NewBookHandler(repo, service)

	body := `{"title":"Go","author":"Alan"}`
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handler.CreateBook(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var resp models.Book
	_ = json.NewDecoder(rec.Body).Decode(&resp)

	if resp.Title != "Go" {
		t.Errorf("unexpected title")
	}
}

func TestUpdateBook(t *testing.T) {
	repo := &mockBookRepo{
		UpdateFn: func(ctx context.Context, book *models.Book) error {
			book.ID = 1
			return nil
		},
	}
	//service := &services.BookService{}
	handler := NewBookHandler(repo, nil)
	body := `{"title":"Go","author":"Alan"}`
	req := httptest.NewRequest(http.MethodPut, "/books/1", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handler.UpdateBook(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp models.Book
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Title != "Go" {
		t.Errorf("unexpected title")
	}
}
