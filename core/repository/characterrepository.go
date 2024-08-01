package repository

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"sanctuary-api/entities"
)

func CheckEquipmentEmplacement(playerEquipments entities.Equipment, emplacement string) bool {
	switch emplacement {
	case "Helmet":
		return playerEquipments.Helmet > 0
	case "Chestplate":
		return playerEquipments.Chestplate > 0
	case "Leggings":
		return playerEquipments.Leggings > 0
	case "Boots":
		return playerEquipments.Boots > 0
	case "Mainhand":
		return playerEquipments.Mainhand > 0
	case "Offhand":
		return playerEquipments.Offhand > 0
	case "AccessorySlot0":
		return playerEquipments.AccessorySlot0 > 0
	case "AccessorySlot1":
		return playerEquipments.AccessorySlot1 > 0
	case "AccessorySlot2":
		return playerEquipments.AccessorySlot2 > 0
	case "AccessorySlot3":
		return playerEquipments.AccessorySlot3 > 0
	default:
		return false
	}
}
func DoUpsertItemInInventory(ctx context.Context, itemID int, playerID int, quantity int, db *pgxpool.Pool, c *gin.Context) {
	_, err := db.Exec(ctx, `
INSERT INTO inventory (item_id, quantity, character_id) VALUES ($1, $2 ,$3)
ON CONFLICT (item_id) DO UPDATE
SET quantity = inventory.quantity + EXCLUDED.quantity;`, itemID, quantity, playerID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request during upsert item into inventory")
		return
	}
}
func DoUpdateEquipment(ctx context.Context, itemID int, playerID int, emplacement string, db *pgxpool.Pool, c *gin.Context) {
	query := fmt.Sprintf("UPDATE equipment SET %s = $1 WHERE player_id = $2", emplacement)
	_, err := db.Exec(ctx, query, itemID, playerID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request during updating equipment"+err.Error())
		return
	}
}

func GetCharactersByID(ctx context.Context, db *pgxpool.Pool, playerID string) (entities.Characters, error) {
	var player entities.Characters
	err := pgxscan.Get(ctx, db, &player, `SELECT id, user_id, email, username, race_id, job_id, exp, po, level, guild_id, inventory_size, location_id FROM characters where id = $1`, playerID)
	if err != nil {
		return player, err
	}

	return player, nil
}

func GetCharactersByEmail(ctx context.Context, db *pgxpool.Pool, playerEmail string) (entities.Characters, error) {
	var player entities.Characters
	err := pgxscan.Get(ctx, db, &player, `SELECT id, user_id, email, username, race_id, job_id, exp, po, level, guild_id, inventory_size, location_id FROM characters where email = $1`, playerEmail)
	if err != nil {
		return player, err
	}

	return player, nil
}

func GetCharactersByLocation(ctx context.Context, db *pgxpool.Pool, locationID string) ([]entities.Characters, error) {
	var players []entities.Characters
	err := pgxscan.Select(ctx, db, &players, `SELECT id, user_id, email, username, race_id, job_id, exp, po, level, guild_id, inventory_size, location_id FROM characters where location_id = $1`, locationID)
	if err != nil {
		return players, err
	}
	return players, nil
}

func UpdateCharactersLocation(ctx context.Context, db *pgxpool.Pool, locationID int, players []entities.Characters) error {
	for _, player := range players {
		_, err := db.Exec(ctx, `UPDATE players SET location_id = $2 where id = $1`, player.ID, locationID)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetCharactersWithItem(ctx context.Context, db *pgxpool.Pool, itemID int) ([]entities.Characters, error) {
	var players []entities.Characters
	err := pgxscan.Select(ctx, db, &players, `SELECT id, user_id, email, username, race_id, job_id, exp, po, level, guild_id, inventory_size, location_id FROM characters c 
    join inventory i on c.id = i.character_id where item_id = $1 UNION
	SELECT id, user_id, email, username, race_id, job_id, exp, po, level, guild_id, inventory_size, location_id  
	FROM characters c JOIN equipment e on c.id = e.character_id
	WHERE e.helmet = $1 OR e.chestplate = $1 OR e.leggings = $1 OR e.boots = $1 OR e.mainhand = $1 OR e.offhand = $1
	OR e.accessory_slot_0 = $1 OR e.accessory_slot_1 = $1 OR e.accessory_slot_2 = $1 OR e.accessory_slot_3 = $1;`, itemID)
	if err != nil {
		return players, err
	}
	fmt.Println(players)
	return players, nil
}

func GetCharactersInventoryByID(ctx context.Context, db *pgxpool.Pool, playerID int) ([]entities.Inventory, error) {
	var inventory []entities.Inventory
	err := pgxscan.Select(ctx, db, &inventory, `SELECT character_id, item_id, quantity FROM inventory where character_id = $1`, playerID)
	if err != nil {
		return inventory, nil
	}
	return inventory, nil
}

func GetCharactersEquipmentByID(ctx context.Context, db *pgxpool.Pool, playerID int) (entities.Equipment, error) {
	var equipment entities.Equipment
	err := pgxscan.Get(ctx, db, &equipment, `SELECT character_id, helmet, chestplate, leggings, boots, mainhand, offhand, accessory_slot_0, accessory_slot_1, accessory_slot_2, accessory_slot_3 FROM equipment where character_id = $1`, playerID)
	if err != nil {
		return equipment, err
	}
	return equipment, nil
}
