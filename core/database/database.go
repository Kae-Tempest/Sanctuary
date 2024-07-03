package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func Connect() *pgxpool.Pool {
	ctx := context.Background()
	db, _ := pgxpool.New(ctx, os.Getenv("DB_URL"))

	return db
}
