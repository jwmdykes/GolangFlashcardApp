package flashcards

import (
	"GolandFlashcardApp/backend/db"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetFlashcard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx := r.Context()
	dbCtx := ctx.Value(db.DbContextKey).(*db.DbContext)

	flashcard, err := dbCtx.GetFlashcard(ctx, id)
	if err != nil {
		http.Error(w, "Flashcard not found", http.StatusNotFound)
		return
	}

	flashcardJSON, _ := json.Marshal(flashcard)
	w.Header().Set("Content-Type", "application/json")
	w.Write(flashcardJSON)
}

func UpdateFlashcard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var dto db.FlashcardDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	ctx := r.Context()
	dbCtx := ctx.Value(db.DbContextKey).(*db.DbContext)

	err = dbCtx.UpdateFlashcard(ctx, id, dto)
	if err != nil {
		http.Error(w, "Flashcard not found", http.StatusNotFound)
		return
	}
}

func DeleteFlashcard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx := r.Context()
	dbCtx := ctx.Value(db.DbContextKey).(*db.DbContext)

	err := dbCtx.DeleteFlashcard(ctx, id)
	if err != nil {
		http.Error(w, "Flashcard not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetFlashcards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dbCtx := ctx.Value(db.DbContextKey).(*db.DbContext)

	flashcards, err := dbCtx.GetFlashcards(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch flashcards", http.StatusInternalServerError)
		return
	}

	flashcardsJSON, _ := json.Marshal(flashcards)
	w.Header().Set("Content-Type", "application/json")
	w.Write(flashcardsJSON)
}

func CreateFlashcard(w http.ResponseWriter, r *http.Request) {
	var dto db.FlashcardDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	dbCtx := ctx.Value(db.DbContextKey).(*db.DbContext)

	createdFlashcard, err := dbCtx.CreateFlashcard(ctx, dto)
	if err != nil {
		http.Error(w, "Failed to create flashcard", http.StatusInternalServerError)
		return
	}

	responseJSON, _ := json.Marshal(createdFlashcard)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
