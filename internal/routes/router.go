package routes

import (
	"net/http"

	"github.com/anastasikalen/bookshelf.git/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(bookHandler *handlers.BookHandler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/books", func(r chi.Router) {
		r.Post("/", bookHandler.CreateBook)
		r.Get("/", bookHandler.GetAllBooks)
		r.Get("/{id}", bookHandler.GetBookByID)
		r.Put("/{id}", bookHandler.UpdateBook)
		r.Delete("/{id}", bookHandler.DeleteBook)

		r.Post("/{id}/mark-out-of-stock", bookHandler.MarkOutOfStock)
		r.Get("/recommend", bookHandler.GetRecommendations)
	})
	return r
}
