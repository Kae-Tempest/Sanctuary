package repository

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"sanctuary-api/entities"
)

func GetLocationByID(ctx context.Context, db *pgxpool.Pool, locationID string) (entities.Locations, error) {
	var location entities.Locations
	err := pgxscan.Get(ctx, db, &location, `SELECT id, type, name, is_safety, size, difficulty FROM locations where id = $1`, locationID)
	if err != nil {
		return location, err
	}
	return location, nil
}

func GetLocationByName(ctx context.Context, db *pgxpool.Pool, locationName string) (entities.Locations, error) {
	var location entities.Locations
	err := pgxscan.Get(ctx, db, &location, `SELECT id, type, name, is_safety, size, difficulty FROM locations where id = $1`, locationName)
	if err != nil {
		return location, err
	}
	return location, nil
}
