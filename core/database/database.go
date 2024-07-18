package database

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() *pgxpool.Pool {
	ctx := context.Background()
	db, err := pgxpool.New(ctx, os.Getenv("DB_URL"))
	if err != nil {
		slog.Error("Error during database connection !")
		panic(err)
	}
	return db
}
