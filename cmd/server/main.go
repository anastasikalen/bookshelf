package main

import (
	"log"
	"net/http"
	"os"

	"github.com/anastasikalen/bookshelf.git/internal/db"
	"github.com/anastasikalen/bookshelf.git/internal/handlers"
	"github.com/anastasikalen/bookshelf.git/internal/repository"
	"github.com/anastasikalen/bookshelf.git/internal/routes"
	"github.com/anastasikalen/bookshelf.git/internal/services"
)

func main() {

	db.ConnectDB()
	defer db.CloseDB()

	db.Migrate()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	repo := repository.NewBookRepository(db.DB)
	service := services.NewBookService(repo)
	handler := handlers.NewBookHandler(repo, service)
	router := routes.NewRouter(handler)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
