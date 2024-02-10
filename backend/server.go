package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"GolandFlashcardApp/backend/db"
	"GolandFlashcardApp/backend/flashcards"
)

func setupRouter(dbCtx *db.DbContext) *chi.Mux {
	var r *chi.Mux = chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(db.DbContextMiddleware(dbCtx)) // Apply the DbContextMiddleware globally

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Route("/flashcards", func(r chi.Router) {
		r.Get("/", flashcards.GetFlashcards)
		r.Post("/", flashcards.CreateFlashcard)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", flashcards.GetFlashcard)
			r.Put("/", flashcards.UpdateFlashcard)
			r.Delete("/", flashcards.DeleteFlashcard)
		})
	})

	return r
}

func main() {
	dbPool, err := pgxpool.New(context.Background(), "postgres://chat_app:chat_app@localhost:5432/chat_app")
	if err != nil {
		panic("Could not connect to the database")
	}
	defer dbPool.Close()

	dbCtx := &db.DbContext{
		DbPool: dbPool,
	}

	r := setupRouter(dbCtx)

	http.ListenAndServe(":5111", r)
}
