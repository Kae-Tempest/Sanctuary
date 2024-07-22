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

func GetLocationByID(ctx context.Context, db *pgxpool.Pool, locationID string) (entities.Locations, error) {
	var location entities.Locations
	err := pgxscan.Get(ctx, db, &location, `SELECT * FROM locations where id = $1`, locationID)
	if err != nil {
		return location, err
	}
	return location, nil
}

func GetLocationByName(ctx context.Context, db *pgxpool.Pool, locationName string) (entities.Locations, error) {
	var location entities.Locations
	err := pgxscan.Get(ctx, db, &location, `SELECT * FROM locations where id = $1`, locationName)
	if err != nil {
		return location, err
	}
	return location, nil
}

func GetSkillByID(ctx context.Context, db *pgxpool.Pool, skillID string) (entities.Skill, error) {
	var skill entities.Skill
	err := pgxscan.Get(ctx, db, &skill, `SELECT * FROM skills where id = $1`, skillID)
	if err != nil {
		return skill, err
	}
	return skill, nil
}

func GetPlayersByLocation(ctx context.Context, db *pgxpool.Pool, locationID string) ([]entities.Player, error) {
	var players []entities.Player
	err := pgxscan.Select(ctx, db, &players, `SELECT * FROM players where location_id = $1`, locationID)
	if err != nil {
		return players, err
	}
	return players, nil
}

func GetCreaturesByLocation(ctx context.Context, db *pgxpool.Pool, locationID string) ([]entities.Creatures, error) {
	var creatures []entities.Creatures
	err := pgxscan.Select(ctx, db, &creatures, `SELECT * FROM creatures c join creaturespawn cp on c.id = cp.creature_id where emplacement_id = $1`, locationID)
	if err != nil {
		return creatures, err
	}
	return creatures, nil
}

func UpdatePlayersLocation(ctx context.Context, db *pgxpool.Pool, locationID int, players []entities.Player) error {
	for _, player := range players {
		_, err := db.Exec(ctx, `UPDATE players SET location_id = $2 where id = $1`, player.ID, locationID)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateCreaturesLocation(ctx context.Context, db *pgxpool.Pool, locationID int, creatures []entities.Creatures) error {
	for _, creature := range creatures {
		_, err := db.Exec(ctx, `UPDATE creaturespawn SET emplacement_id = $2 where creature_id = $1`, creature.ID, locationID)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetItemByID(ctx context.Context, db *pgxpool.Pool, itemID string) (entities.Items, error) {
	var item entities.Items
	err := pgxscan.Get(ctx, db, &item, `SELECT * FROM items where id = $1`, itemID)
	if err != nil {
		return item, nil
	}

	return item, nil
}

func GetItemByName(ctx context.Context, db *pgxpool.Pool, itemName string) (entities.Items, error) {
	var item entities.Items
	err := pgxscan.Get(ctx, db, &item, `SELECT * FROM items where name = $1`, itemName)
	if err != nil {
		return item, nil
	}

	return item, nil
}

func GetItemsByType(ctx context.Context, db *pgxpool.Pool, itemType string) ([]entities.Items, error) {
	var items []entities.Items
	err := pgxscan.Select(ctx, db, &items, `SELECT * FROM items where type = $1`, itemType)
	if err != nil {
		return items, nil
	}

	return items, nil
}

func GetPlayersWithItem(ctx context.Context, db *pgxpool.Pool, itemName string) ([]entities.Player, error) {
	var players []entities.Player
	return players, nil
}

func GetPlayerInventoryByID(ctx context.Context, db *pgxpool.Pool, playerID int) ([]entities.Inventory, error) {
	var inventory []entities.Inventory
	return inventory, nil
}

func GetPlayerEquipmentByID(ctx context.Context, db *pgxpool.Pool, playerID int) (entities.Equipment, error) {
	var equipment entities.Equipment
	return equipment, nil
}
