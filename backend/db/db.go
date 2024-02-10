package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DbContext struct {
	DbPool *pgxpool.Pool
}

type Flashcard struct {
	Id          string    `json:"id"`
	Question    string    `json:"question"`
	Answer      string    `json:"answer"`
	DateCreated time.Time `json:"date_created"`
}

type FlashcardDTO struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func (db *DbContext) UpdateFlashcard(ctx context.Context, id string, dto FlashcardDTO) error {
	sql :=
		`
	UPDATE flashcards SET 
	question = CASE WHEN $1 = '' THEN question ELSE $1 END,
	answer = CASE WHEN $2 = '' THEN answer ELSE $2 END
	WHERE id = $3
	`

	cmdTag, err := db.DbPool.Exec(ctx, sql, dto.Question, dto.Answer, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *DbContext) DeleteFlashcard(ctx context.Context, id string) error {
	sql := `DELETE FROM flashcards WHERE id = $1`

	cmdTag, err := db.DbPool.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *DbContext) GetFlashcard(ctx context.Context, id string) (*Flashcard, error) {
	var flashcard Flashcard
	err := db.DbPool.QueryRow(ctx, `SELECT id, question, answer FROM flashcards WHERE id = $1`, id).Scan(&flashcard.Id, &flashcard.Question, &flashcard.Answer)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &flashcard, nil
}

func (db *DbContext) GetFlashcards(ctx context.Context) ([]Flashcard, error) {
	flashcards := []Flashcard{}

	rows, err := db.DbPool.Query(ctx, `SELECT id, question, answer FROM flashcards`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f Flashcard
		if err := rows.Scan(&f.Id, &f.Question, &f.Answer); err != nil {
			return nil, err
		}
		flashcards = append(flashcards, f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return flashcards, nil
}

func (db *DbContext) CreateFlashcard(ctx context.Context, dto FlashcardDTO) (*Flashcard, error) {
	sql := `INSERT INTO flashcards (question, answer) VALUES ($1, $2) RETURNING id, question, answer`

	var createdFlashcard Flashcard

	err := db.DbPool.QueryRow(ctx, sql, dto.Question, dto.Answer).Scan(&createdFlashcard.Id, &createdFlashcard.Question, &createdFlashcard.Answer)
	if err != nil {
		return nil, err
	}

	return &createdFlashcard, nil
}
