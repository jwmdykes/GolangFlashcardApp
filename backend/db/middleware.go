package db

import (
	"context"
	"net/http"
)

type contextKey string

const DbContextKey contextKey = "dbContext"

func DbContextMiddleware(db *DbContext) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), DbContextKey, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
