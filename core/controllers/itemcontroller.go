package controllers

import (
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"net/http"
	"sanctuary-api/database"
	"sanctuary-api/entities"
	"sanctuary-api/repository"
	"strconv"
)

func GetItems(c *gin.Context) {
	db := database.Connect()

	var items []entities.ItemComplete
	err := pgxscan.Select(ctx, db, &items, `SELECT * FROM items join item_stats on items.id = item_stats.item_id`)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
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
	var itemForm entities.ItemComplete
	if err := c.ShouldBindBodyWithJSON(&itemForm); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	_, err := repository.GetItemByName(ctx, db, itemForm.Item.Name)
	if errors.Is(pgx.ErrNoRows, err) {
		_, insertErr := db.Exec(ctx, `INSERT INTO items (name, description, type) values ($1,$2,$3)`,
			itemForm.Item.Name, itemForm.Item.Description, itemForm.Item.Type)
		if insertErr != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}

		item, getErr := repository.GetItemByName(ctx, db, itemForm.Item.Name)
		if getErr != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}

		_, insertErr = db.Exec(ctx, `INSERT INTO item_stats (item_id, strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma, enchantment_level) VALUES  ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
			item.ID, itemForm.Stats.Strength, itemForm.Stats.Constitution, itemForm.Stats.Mana, itemForm.Stats.Stamina, itemForm.Stats.Dexterity,
			itemForm.Stats.Intelligence, itemForm.Stats.Wisdom, itemForm.Stats.Charisma, itemForm.Stats.EnchantmentLevel)

		completeItem, cErr := repository.GetCompleteItem(ctx, db, item.ID)
		if cErr != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}

		c.JSON(http.StatusCreated, &completeItem)
	}
	if err != nil && !errors.Is(pgx.ErrNoRows, err) {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusConflict, "already exist")
}

func UpdateItem(c *gin.Context) {
	db := database.Connect()
	var itemForm entities.Items
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
    = ($2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) where id = $1`, item.Item.ID, itemForm.Name, itemForm.Description, itemForm.Type)
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
    = ($2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) where item_id = $1`, item.Item.ID, itemForm.Strength, itemForm.Constitution, itemForm.Mana, itemForm.Stamina, itemForm.Dexterity,
		itemForm.Intelligence, itemForm.Wisdom, itemForm.Charisma, itemForm.EnchantmentLevel)
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	updatedItem, getErr := repository.GetItemStat(ctx, db, item.Item.ID)
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

	// get Player with item
	players, getErr := repository.GetPlayersWithItem(ctx, db, item.Item.ID)
	if getErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	for _, player := range players {
		// get inventory
		inventory, invErr := repository.GetPlayerInventoryByID(ctx, db, player.ID)
		if invErr != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		// del item in inventory
		for _, invItem := range inventory {
			_, delErr := db.Exec(ctx, `DELETE FROM inventory where player_id = $1 AND item_id = $2`, player.ID, invItem.ItemID)
			if delErr != nil {
				c.JSON(http.StatusBadRequest, "bad request")
				return
			}
		}
		// get equipment
		equipment, equipErr := repository.GetPlayerEquipmentByID(ctx, db, player.ID)
		if equipErr != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}

		switch {
		case equipment.Helmet == item.Item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "helmet", db, c)
		case equipment.Chestplate == item.Item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "chestplate", db, c)
		case equipment.Leggings == item.Item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "leggings", db, c)
		case equipment.Boots == item.Item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "boots", db, c)
		case equipment.Mainhand == item.Item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "mainhand", db, c)
		case equipment.Offhand == item.Item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "offhand", db, c)
		case equipment.AccessorySlot0 == item.Item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "accessory_slot_0", db, c)
		case equipment.AccessorySlot1 == item.Item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "accessory_slot_1", db, c)
		case equipment.AccessorySlot2 == item.Item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "accessory_slot_2", db, c)
		case equipment.AccessorySlot3 == item.Item.ID:
			repository.DoUpdateEquipment(ctx, 0, player.ID, "accessory_slot_3", db, c)
		}
	}

	_, deleteErr := db.Exec(ctx, `DELETE FROM item_stats where item_id = $1`, item.Item.ID)
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	_, deleteErr = db.Exec(ctx, `DELETE FROM items where id = $1`, item.Item.ID)
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.Status(http.StatusOK)

}
