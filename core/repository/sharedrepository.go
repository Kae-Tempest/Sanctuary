package repository

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"sanctuary-api/entities"
)

func GetPlayerByID(ctx context.Context, db *pgxpool.Pool, playerID string) (entities.Player, error) {
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

func GetCreatureByID(ctx context.Context, db *pgxpool.Pool, creatureID string) (entities.Creatures, error) {
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

func GetLocationByID(ctx context.Context, db *pgxpool.Pool, locationID int) (entities.Locations, error) {
	var location entities.Locations
	err := pgxscan.Get(ctx, db, &location, `SELECT * FROM locations where id = $1`, locationID)
	if err != nil {
		return location, err
	}
	return location, nil
}

func GetSkillByID(ctx context.Context, db *pgxpool.Pool, skillID int) (entities.Skill, error) {
	var skill entities.Skill
	err := pgxscan.Get(ctx, db, &skill, `SELECT * FROM skills where id = $1`, skillID)
	if err != nil {
		return skill, err
	}
	return skill, nil
}
