package repository

import (
	"context"
	"fmt"
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
INSERT INTO inventory (item_id, quantity, player_id) VALUES ($1, $2 ,$3)
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
