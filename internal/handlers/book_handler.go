package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/anastasikalen/bookshelf.git/internal/models"
	"github.com/anastasikalen/bookshelf.git/internal/repository"
	"github.com/anastasikalen/bookshelf.git/internal/services"
	"github.com/go-chi/chi/v5"
)

type BookHandler struct {
	repo repository.BookRepository
	serv *services.BookService
}

func NewBookHandler(
	repo repository.BookRepository,
	service *services.BookService,
) *BookHandler {
	return &BookHandler{
		repo: repo,
		serv: service,
	}
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if book.Title == "" || book.Author == "" {
		log.Printf("title and author are required")
		http.Error(w, "title and author are required", http.StatusBadRequest)
		return
	}

	if len(book.Title) > 500 {
		log.Printf("title is too long")
		http.Error(w, "title must be less than 500 characters", http.StatusBadRequest)
		return
	}

	if len(book.Author) > 200 {
		log.Printf("author is too long")
		http.Error(w, "author must be less than 200 characters", http.StatusBadRequest)
		return
	}

	if book.ISBN != "" && len(book.ISBN) > 20 {
		log.Printf("ISBN is too long")
		http.Error(w, "ISBN must be less than 20 characters", http.StatusBadRequest)
		return
	}

	if book.Rating != nil {
		if *book.Rating < 1 || *book.Rating > 10 {
			http.Error(w, "rating must be between 1 and 10", http.StatusBadRequest)
			return
		}
	}

	if book.Year < 0 || book.Year > 9999 {
		log.Printf("year must be between 0 and 9999, got: %d", book.Year)
		http.Error(w, "year must be between 0 and 9999", http.StatusBadRequest)
		return
	}
	if err := h.repo.Create(r.Context(), &book); err != nil {
		log.Printf("Failed to create book: %v", err)
		http.Error(w, "failed to create book", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		log.Printf("Invalid book ID %v", idParam)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	book, err := h.repo.GetByID(r.Context(), uint(id))
	if err != nil {
		log.Printf("Failed to get book: %v", err)
		http.Error(w, "failed to get book by id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	books, err := h.repo.GetAll(r.Context(), limit, offset)
	if err != nil {
		log.Printf("Failed to get books: %v", err)
		http.Error(w, "failed to get books", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		log.Printf("Invalid book ID %v", idParam)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if book.Title == "" || book.Author == "" {
		log.Printf("title and author are required")
		http.Error(w, "title and author are required", http.StatusBadRequest)
		return
	}

	if len(book.Title) > 500 {
		log.Printf("title is too long")
		http.Error(w, "title must be less than 500 characters", http.StatusBadRequest)
		return
	}

	if len(book.Author) > 200 {
		log.Printf("author is too long")
		http.Error(w, "author must be less than 200 characters", http.StatusBadRequest)
		return
	}

	if book.ISBN != "" && len(book.ISBN) > 20 {
		log.Printf("ISBN is too long")
		http.Error(w, "ISBN must be less than 20 characters", http.StatusBadRequest)
		return
	}

	if book.Rating != nil {
		if *book.Rating < 1 || *book.Rating > 10 {
			http.Error(w, "rating must be between 1 and 10", http.StatusBadRequest)
			return
		}
	}

	if book.Year < 0 || book.Year > 9999 {
		log.Printf("year must be between 0 and 9999, got: %d", book.Year)
		http.Error(w, "year must be between 0 and 9999", http.StatusBadRequest)
		return
	}
	book.ID = uint(id)
	if err := h.repo.Update(r.Context(), &book); err != nil {
		log.Printf("Failed to update book: %v", err)
		http.Error(w, "failed to update book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		log.Printf("Invalid book ID %v", idParam)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.repo.Delete(r.Context(), uint(id)); err != nil {
		log.Printf("Failed to delete book: %v", err)
		http.Error(w, "failed to delete book", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *BookHandler) MarkOutOfStock(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		log.Printf("Invalid book ID %v", idParam)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	err = h.serv.MarkOutOfStock(r.Context(), uint(id))
	if err != nil {
		if err.Error() == "book not found" {
			log.Printf("Failed to mark out of stock: %v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Failed to mark out of stock: %v", err)
		http.Error(w, "failed to update book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BookHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	books, err := h.repo.GetRecommendations(r.Context())
	if err != nil {
		log.Printf("Failed to get recommendations: %v", err)
		http.Error(w, "failed to get recommendations", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(books)

}
