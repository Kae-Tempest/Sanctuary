package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"net/http"
	"sanctuary-api/database"
	"sanctuary-api/entities"
	"sanctuary-api/repository"
	"strconv"
)

func GetItems(c *gin.Context) {
	db := database.Connect()

	rows, queryErr := db.Query(ctx, `SELECT id ,name, description, type, rank,
	    strength,constitution,mana,stamina,dexterity,intelligence,wisdom,charisma,enchantment_level,emplacement
		FROM items full join item_stats  on items.id = item_stats.item_id  full join item_emplacement on items.id = item_emplacement.item_id`)
	if queryErr != nil {
		fmt.Println(queryErr)
	}

	items, err := repository.AssignMultipleRowsItem(rows)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, &items)
}

func GetItemByID(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	item, err := repository.GetItemByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &item)
}

func GetItemByType(c *gin.Context) {
	db := database.Connect()
	itemType := c.Param("type")

	items, err := repository.GetItemsByType(ctx, db, itemType)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &items)
}

func CreateItem(c *gin.Context) {
	db := database.Connect()
	var itemForm entities.Item
	if err := c.ShouldBindBodyWithJSON(&itemForm); err != nil {
		slog.Error("Error during binding ItemForm", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	existingItem, err := repository.GetItemByName(ctx, db, itemForm.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		_, insertErr := db.Exec(ctx, `INSERT INTO items (name, description, type, rank) values ($1,$2,$3,$4)`,
			itemForm.Name, itemForm.Description, itemForm.Type, itemForm.Rank)
		if insertErr != nil {
			slog.Error("Error Inserting Item", slog.Any("error", insertErr))
			c.JSON(http.StatusBadRequest, "Error Inserting Item")
		}

		item, getErr := repository.GetItemInfoByName(ctx, db, itemForm.Name)
		if getErr != nil {
			slog.Error("Error Getting Item by Name", slog.Any("error", getErr))
			c.JSON(http.StatusBadRequest, "Error Getting Item by Name")
			return
		}

		_, insertErr = db.Exec(ctx, `INSERT INTO item_stats (item_id, strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma, enchantment_level) VALUES  ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
			item.ID, itemForm.Stats.Strength, itemForm.Stats.Constitution, itemForm.Stats.Mana, itemForm.Stats.Stamina, itemForm.Stats.Dexterity,
			itemForm.Stats.Intelligence, itemForm.Stats.Wisdom, itemForm.Stats.Charisma, itemForm.Stats.EnchantmentLevel)
		if insertErr != nil {
			slog.Error("Error Inserting Item Stat", slog.Any("error", insertErr))
			c.JSON(http.StatusBadRequest, "Error Inserting Item")
			return
		}

		_, insertErr = db.Exec(ctx, `INSERT INTO item_emplacement (item_id, emplacement) VALUES ($1,$2)`, item.ID, itemForm.ItemEmplacement.Emplacement)
		if insertErr != nil {
			slog.Error("Error Inserting Item Emplacement", slog.Any("error", insertErr))
			c.JSON(http.StatusBadRequest, "Error Inserting Item")
			return
		}

		completeItem, cErr := repository.GetItemByID(ctx, db, strconv.Itoa(item.ID))
		if cErr != nil {
			slog.Error("Error  Getting Item by ID", slog.Any("error", cErr))
			c.JSON(http.StatusBadRequest, "Error Getting Item by ID")
			return
		}

		c.JSON(http.StatusCreated, &completeItem)
		return
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("Error Getting Item by Name", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error Getting Item by Name")
		return
	}
	if existingItem.ID != 0 {
		c.String(http.StatusConflict, "already exist")
		return
	}
}

func UpdateItem(c *gin.Context) {
	db := database.Connect()
	var itemForm entities.Item
	if err := c.ShouldBindBodyWithJSON(&itemForm); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	item, err := repository.GetItemByID(ctx, db, strconv.Itoa(itemForm.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	_, updateErr := db.Exec(ctx, `UPDATE items set (name, description, type)
    = ($2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) where id = $1`, item.ID, itemForm.Name, itemForm.Description, itemForm.Type)
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	updatedItem, gerErr := repository.GetItemByID(ctx, db, strconv.Itoa(itemForm.ID))
	if gerErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &updatedItem)

}

func UpdateItemStat(c *gin.Context) {
	db := database.Connect()
	var itemForm entities.ItemStat
	if err := c.ShouldBindBodyWithJSON(&itemForm); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	item, err := repository.GetItemByID(ctx, db, strconv.Itoa(itemForm.ItemID))
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	_, updateErr := db.Exec(ctx, `UPDATE item_stats set (strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma, enchantment_level)
    = ($2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) where item_id = $1`, item.ID, itemForm.Strength, itemForm.Constitution, itemForm.Mana, itemForm.Stamina, itemForm.Dexterity,
		itemForm.Intelligence, itemForm.Wisdom, itemForm.Charisma, itemForm.EnchantmentLevel)
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	updatedItem, getErr := repository.GetItemStat(ctx, db, item.ID)
	if getErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &updatedItem)
}

func UpdateItemEmplacement(c *gin.Context) {
	db := database.Connect()
	var itemForm entities.ItemEmplacement
	if err := c.ShouldBindBodyWithJSON(&itemForm); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	item, err := repository.GetItemByID(ctx, db, strconv.Itoa(itemForm.ItemID))
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	_, updateErr := db.Exec(ctx, `INSERT INTO item_emplacement (item_id, emplacement) VALUES ($1,$2) on conflict (item_id) do UPDATE SET emplacement = $2`, item.ID, itemForm.Emplacement)
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	updatedItem, getErr := repository.GetItemEmplacement(ctx, db, item.ID)
	if getErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &updatedItem)
}

func DeleteItem(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	item, err := repository.GetItemByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// get Characters with item
	players, getErr := repository.GetCharactersWithItem(ctx, db, item.ID)
	if getErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	for _, player := range players {
		// get inventory
		inventory, invErr := repository.GetCharactersInventoryByID(ctx, db, player.ID)
		if invErr != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		// del item in inventory
		for _, invItem := range inventory {
			_, delErr := db.Exec(ctx, `DELETE FROM inventory where character_id = $1 AND item_id = $2`, player.ID, invItem.ItemID)
			if delErr != nil {
				c.JSON(http.StatusBadRequest, "bad request")
				return
			}
		}
		// get equipment
		equipment, equipErr := repository.GetCharactersEquipmentByID(ctx, db, player.ID)
		if equipErr != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}

		switch {
		case equipment.Helmet == item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "helmet", db, c)
		case equipment.Chestplate == item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "chestplate", db, c)
		case equipment.Leggings == item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "leggings", db, c)
		case equipment.Boots == item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "boots", db, c)
		case equipment.Mainhand == item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "mainhand", db, c)
		case equipment.Offhand == item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "offhand", db, c)
		case equipment.AccessorySlot0 == item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "accessory_slot_0", db, c)
		case equipment.AccessorySlot1 == item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "accessory_slot_1", db, c)
		case equipment.AccessorySlot2 == item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "accessory_slot_2", db, c)
		case equipment.AccessorySlot3 == item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "accessory_slot_3", db, c)
		}
	}

	_, deleteErr := db.Exec(ctx, `DELETE FROM item_stats where item_id = $1`, item.ID)
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	_, deleteErr = db.Exec(ctx, `DELETE FROM item_emplacement where item_id = $1`, item.ID)
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	_, deleteErr = db.Exec(ctx, `DELETE FROM items where id = $1`, item.ID)
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.Status(http.StatusOK)

}
