package repository

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"sanctuary-api/entities"
)

func GetPetByID(ctx context.Context, db *pgxpool.Pool, petID string) (entities.PetsMounts, error) {
	var pet entities.PetsMounts
	err := pgxscan.Get(ctx, db, &pet, `SELECT mob_id, is_mountable, speed, id FROM pets_mounts where id = $1`, petID)
	if err != nil {
		return pet, err
	}
	return pet, nil
}
