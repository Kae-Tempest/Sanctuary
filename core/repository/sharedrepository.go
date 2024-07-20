package repository

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"sanctuary-api/entities"
)

func GetPlayerById(ctx context.Context, db *pgxpool.Pool, playerID string) (entities.Player, error) {
	var player entities.Player
	err := pgxscan.Get(ctx, db, &player, `SELECT * FROM players where id = $1`, playerID)
	if err != nil {
		return player, err
	}

	return player, nil
}

func GetPlayerByEmail(ctx context.Context, db *pgxpool.Pool, playerEmail string) (entities.Player, error) {
	var player entities.Player
	err := pgxscan.Get(ctx, db, &player, `SELECT * FROM players where email = $1`, playerEmail)
	if err != nil {
		return player, err
	}

	return player, nil
}

func GetCreatureById(ctx context.Context, db *pgxpool.Pool, creatureID string) (entities.Creatures, error) {
	var creature entities.Creatures
	err := pgxscan.Get(ctx, db, &creature, `SELECT * FROM creatures where id = $1`, creatureID)
	if err != nil {
		return creature, err
	}
	return creature, nil
}

func GetCreatureByName(ctx context.Context, db *pgxpool.Pool, creatureName string) (entities.Creatures, error) {
	var creature entities.Creatures
	err := pgxscan.Get(ctx, db, &creature, `SELECT * FROM creatures where name = $1`, creatureName)
	if err != nil {
		return creature, err
	}
	return creature, nil
}
