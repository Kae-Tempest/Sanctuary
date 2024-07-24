package repository

import (
	"context"
	"fmt"
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

func GetCreatureByID(ctx context.Context, db *pgxpool.Pool, mobID string) (entities.Mob, error) {
	var mob entities.Mob
	err := pgxscan.Get(ctx, db, &mob, `SELECT * FROM mobs where id = $1`, mobID)
	if err != nil {
		return mob, err
	}
	return mob, nil
}

func GetCreatureByName(ctx context.Context, db *pgxpool.Pool, mobName string) (entities.Mob, error) {
	var mob entities.Mob
	err := pgxscan.Get(ctx, db, &mob, `SELECT * FROM mobs where name = $1`, mobName)
	if err != nil {
		return mob, err
	}
	return mob, nil
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

func GetCreaturesByLocation(ctx context.Context, db *pgxpool.Pool, locationID string) ([]entities.Mob, error) {
	var mobs []entities.Mob
	err := pgxscan.Select(ctx, db, &mobs, `SELECT * FROM mobs c join mob_spawn cp on c.id = cp.mob_id where location_id = $1`, locationID)
	if err != nil {
		return mobs, err
	}
	return mobs, nil
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

func UpdateCreaturesLocation(ctx context.Context, db *pgxpool.Pool, locationID int, mobs []entities.Mob) error {
	for _, mob := range mobs {
		_, err := db.Exec(ctx, `UPDATE creaturespawn SET emplacement_id = $2 where creature_id = $1`, mob.ID, locationID)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetItemByID(ctx context.Context, db *pgxpool.Pool, itemID string) (entities.ItemComplete, error) {
	var item entities.ItemComplete
	err := pgxscan.Get(ctx, db, &item, `SELECT * FROM items join item_stats on items.id = item_stats.item_id where id = $1`, itemID)
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

func GetCompleteItem(ctx context.Context, db *pgxpool.Pool, itemID int) (entities.ItemComplete, error) {
	var item entities.ItemComplete
	err := pgxscan.Get(ctx, db, &item, `SELECT * FROM items i join item_stats is on i.id = is.item_id where id = $1`, itemID)
	if err != nil {
		return item, nil
	}

	return item, nil
}

func GetItemStat(ctx context.Context, db *pgxpool.Pool, itemID int) (entities.ItemStat, error) {
	var item entities.ItemStat
	err := pgxscan.Get(ctx, db, &item, `SELECT * FROM item_stats where item_id = $1`, itemID)
	if err != nil {
		return item, nil
	}

	return item, nil
}

func GetPlayersWithItem(ctx context.Context, db *pgxpool.Pool, itemID int) ([]entities.Player, error) {
	var players []entities.Player
	err := pgxscan.Select(ctx, db, &players, `SELECT p.*FROM players p join inventory i on p.id = i.player_id where item_id = $1 UNION
SELECT p.*  FROM players p JOIN equipment e on p.id = e.player_id
WHERE e.helmet = $1 OR e.chestplate = $1 OR e.leggings = $1 OR e.boots = $1 OR e.mainhand = $1 OR e.offhand = $1
OR e.accessory_slot_0 = $1 OR e.accessory_slot_1 = $1 OR e.accessory_slot_2 = $1 OR e.accessory_slot_3 = $1;`, itemID)
	if err != nil {
		return players, err
	}
	fmt.Println(players)
	return players, nil
}

func GetPlayerInventoryByID(ctx context.Context, db *pgxpool.Pool, playerID int) ([]entities.Inventory, error) {
	var inventory []entities.Inventory
	err := pgxscan.Select(ctx, db, &inventory, `SELECT * FROM inventory where player_id = $1`, playerID)
	if err != nil {
		return inventory, nil
	}
	return inventory, nil
}

func GetPlayerEquipmentByID(ctx context.Context, db *pgxpool.Pool, playerID int) (entities.Equipment, error) {
	var equipment entities.Equipment
	err := pgxscan.Get(ctx, db, &equipment, `SELECT * FROM equipment where player_id = $1`, playerID)
	if err != nil {
		return equipment, err
	}
	return equipment, nil
}
