package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sanctuary-api/database"
	"sanctuary-api/entities"
	"sanctuary-api/repository"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

var ctx = context.Background()

// GET \\
func GetAllCharacters(c *gin.Context) {
	db := database.Connect()
	var players []entities.Characters
	err := pgxscan.Select(ctx, db, &players, `SELECT id, email, username, race_id, job_id, exp, level, guild_id, inventory_size, po, location_id, user_id FROM characters`)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	if len(players) > 0 {
		c.JSON(http.StatusOK, &players)
	} else {
		c.JSON(http.StatusNotFound, gin.H{})
	}
}
func GetOneCharacters(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	player, err := repository.GetCharactersByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &player)
}
func GetCharactersStats(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	player, err := repository.GetCharactersByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var playerStats entities.Stats
	err = pgxscan.Get(ctx, db, &playerStats, `SELECT character_id, strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma, hp from stats where character_id = $1`, player.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerStats)

}
func GetCharactersEquipment(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	player, err := repository.GetCharactersByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var playerEquipment entities.Equipment
	err = pgxscan.Get(ctx, db, &playerEquipment, `SELECT character_id, helmet, chestplate, leggings, boots, mainhand, offhand, accessory_slot_0, accessory_slot_1, accessory_slot_2, accessory_slot_3 from equipment where character_id = $1`, player.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerEquipment)

}
func GetCharactersInventory(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	player, err := repository.GetCharactersByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var playerInventory entities.Inventory
	err = pgxscan.Get(ctx, db, &playerInventory, `SELECT character_id, item_id, quantity from inventory where character_id = $1`, player.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerInventory)

}
func GetCharactersPets(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	player, err := repository.GetCharactersByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var playerPets []entities.CharactersPet
	err = pgxscan.Select(ctx, db, &playerPets, `SELECT pet_id FROM character_pets_mounts where character_id = $1`, player.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerPets)

}
func GetCharactersGuild(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	player, err := repository.GetCharactersByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var playerGuild entities.Guild
	err = pgxscan.Get(ctx, db, &playerGuild, `SELECT g.id, g.name, g.owner, array_agg(gm.character_id)  FROM guilds g join guilds_members gm on g.id = gm.guilds_id where gm.character_id = $1 group by g.owner, g.name, g.id`, player.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerGuild)

}
func GetCharactersSkill(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	player, err := repository.GetCharactersByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var playerSkill []entities.Skill
	err = pgxscan.Select(ctx, db, &playerSkill, `SELECT character_id, skill_id from character_skill where character_id = $1`, player.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerSkill)

}

// POST \\
func CreateCharacters(c *gin.Context) {
	db := database.Connect()
	var playerForm entities.Characters
	if err := c.ShouldBindBodyWithJSON(&playerForm); err != nil {
		c.String(http.StatusBadRequest, "bad request during bind JSON")
		return
	}

	existingCharacters, epErr := repository.GetCharactersByEmail(ctx, db, playerForm.Email)
	if errors.Is(pgx.ErrNoRows, epErr) {
		c.JSON(http.StatusConflict, "existing player"+epErr.Error())
		fmt.Println(existingCharacters)
	}
	if epErr != nil && !errors.Is(epErr, pgx.ErrNoRows) {
		c.String(http.StatusBadRequest, "bad request")
	}

	_, err := db.Exec(ctx, `INSERT INTO characters (email, username, race_id, job_id, exp, level, guild_id, inventory_size, po, location_id, user_id) values ($1, $2, $3, $4, 0, 1, 0, 10, 50, 1, null)`,
		playerForm.Email, playerForm.Username, playerForm.RaceID, playerForm.JobID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request during inserting data"+err.Error())
		return
	}
	var playerID int
	err = pgxscan.Get(ctx, db, &playerID, `SELECT id from characters where email = $1`, playerForm.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request during get CharactersID"+err.Error())
		return
	}
	// Get Job Stats
	var playerJob entities.Job
	err = pgxscan.Get(ctx, db, &playerJob, `SELECT id, name, description, strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma from jobs where id = $1`, playerForm.JobID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request during get Job"+err.Error())
		return
	}
	// Get Race Stats
	var playerRace entities.Race
	err = pgxscan.Get(ctx, db, &playerRace, `SELECT id, name, description, mana, stamina, wisdom, charisma from races where id = $1`, playerForm.RaceID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request during get Race"+err.Error())
		return
	}
	// Create Characters Stats
	_, err = db.Exec(ctx, `INSERT INTO stats (character_id, strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma, hp) values ($1, $2, $3, $4,$5,$6,$7,$8,$9,20)`,
		playerID, playerJob.Strength, playerJob.Constitution, playerRace.Mana, playerRace.Stamina, playerJob.Dexterity, playerJob.Intelligence, playerRace.Wisdom, playerRace.Charisma)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request during insert Stat"+err.Error())
		return
	}

	_, err = db.Exec(ctx, `INSERT INTO equipment (character_id, helmet, chestplate, leggings, boots, mainhand, offhand, accessory_slot_0, accessory_slot_1, accessory_slot_2, accessory_slot_3) VALUES 
    ($1,0,0,0,0,0,0,0,0,0,0)`, playerID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request during insert equipment"+err.Error())
		return
	}

	c.Status(http.StatusCreated)
}
func AddItemToCharactersInventory(c *gin.Context) {
	type Body struct {
		ItemID   int
		Quantity int
	}
	db := database.Connect()
	id := c.Param("id")
	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// verifier si l'item existe
	var selectedItem entities.Item
	err := pgxscan.Get(ctx, db, &selectedItem, `SELECT id, name, description, type, rank FROM items where id = $1`, body.ItemID)
	if err != nil {
		c.JSON(http.StatusNotFound, "Item with this ID doesn't exist !")
		return
	}

	player, err := repository.GetCharactersByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// verif if user had already item
	var playerInventory []entities.Inventory
	err = pgxscan.Select(ctx, db, &playerInventory, `SELECT item_id, quantity FROM inventory where character_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	if len(playerInventory) <= 0 {
		_, err = db.Exec(ctx, `INSERT into inventory (character_id, item_id, quantity) values ($1,$2,$3)`, player.ID, body.ItemID, body.Quantity)
		if err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		c.JSON(http.StatusCreated, "First item add to inventory")
	}
	var haveItem = false
	var quantity int
	for _, item := range playerInventory {
		if item.ItemID == body.ItemID {
			quantity = item.Quantity
			haveItem = true
			break
		} else {
			haveItem = false
			continue
		}
	}

	if haveItem {
		_, err = db.Exec(ctx, `UPDATE inventory SET quantity = $1 where player_id = $2`, quantity+body.Quantity, player.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		c.JSON(http.StatusCreated, "Item quantity updated")
	} else {
		_, err = db.Exec(ctx, `INSERT into inventory (character_id, item_id, quantity) values ($1,$2,$3)`, player.ID, body.ItemID, body.Quantity)
		if err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		c.JSON(http.StatusCreated, "Item add to inventory")
	}
}
func AddPetToCharacters(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	type Body struct {
		PetsID int
	}

	var playerPetForm Body
	err := c.ShouldBindBodyWithJSON(&playerPetForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// check if player exist
	player, err := repository.GetCharactersByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	// check if pet exist
	var selectedPet entities.PetsMounts
	err = pgxscan.Get(ctx, db, &selectedPet, `SELECT id FROM pets_mounts where id = $1`, playerPetForm.PetsID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request selecting pets")
		return
	}
	// check if player already have this pet
	var playerPets entities.CharactersPet
	err = pgxscan.Get(ctx, db, &playerPets, `SELECT pet_id FROM character_pets_mounts where pet_id = $1`, playerPetForm.PetsID)
	if errors.Is(err, pgx.ErrNoRows) {
		_, insertErr := db.Exec(ctx, `INSERT INTO character_pets_mounts (pet_id, character_id) values ($1,$2)`, playerPetForm.PetsID, player.ID)
		if insertErr != nil {
			c.String(http.StatusBadRequest, "bad request inserting user's pets")
			return
		}
		c.JSON(http.StatusOK, "New Pet Added")
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		c.String(http.StatusBadRequest, "bad request selecting user's pets")
		return
	}
	if playerPets.PetID != 0 {
		var mob entities.Mob
		err = pgxscan.Get(ctx, db, &mob, `SELECT name from mobs where id = $1`, selectedPet.MobID)
		if err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}

		var petScroll entities.Item
		err = pgxscan.Get(ctx, db, &petScroll, `SELECT id FROM items where name like '%$1%'`, mob.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}

		repository.DoUpsertItemInInventory(ctx, petScroll.ID, player.ID, 1, db, c)

		c.JSON(http.StatusOK, "You Already have this pet, you got his scroll")
	}
}
func AddSkillToCharacters(c *gin.Context) {
	type Body struct {
		SkillId int
	}
	db := database.Connect()
	id := c.Param("id")
	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var selectedSkill entities.Skill
	err := pgxscan.Get(ctx, db, &selectedSkill, `SELECT id FROM skills where id = $1`, body.SkillId)
	if err != nil {
		c.JSON(http.StatusNotFound, "Pet with this ID doesn't exist !")
		return
	}

	player, err := repository.GetCharactersByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var playerSkills []entities.CharactersSkill
	err = pgxscan.Select(ctx, db, &playerSkills, `SELECT skill_id FROM character_skill where character_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	if len(playerSkills) <= 0 {
		_, err = db.Exec(ctx, `INSERT into character_skill (character_id, skill_id) values ($1,$2)`, player.ID, body.SkillId)
		if err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		c.JSON(http.StatusCreated, "First pet get by User")
	}

	var haveSkill = false
	for _, skill := range playerSkills {
		if skill.SkillID == body.SkillId {
			haveSkill = true
			break
		} else {
			haveSkill = false
			continue
		}
	}

	if haveSkill {
		c.JSON(http.StatusBadRequest, "User already have this skill")
	} else {
		_, err = db.Exec(ctx, `INSERT into character_skill (character_id, skill_id) values ($1,$2)`, player.ID, body.SkillId)
		if err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		c.JSON(http.StatusCreated, "User tamed this skill")
	}

}

// PATCH \\
func UpdateCharactersStats(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	var playerStatsForm entities.Stats
	if err := c.ShouldBindBodyWithJSON(&playerStatsForm); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	player, err := repository.GetCharactersByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var playerStats entities.Stats
	err = pgxscan.Get(ctx, db, &playerStats, `SELECT character_id, strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma, hp FROM stats where character_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	_, err = db.Exec(ctx, `UPDATE stats set (strength , charisma , constitution, dexterity, hp, intelligence, mana, stamina, wisdom) = ($2,$3,$4,$5,$6,$7,$8,$9,$10) where player_id = $1`,
		id, playerStats.Strength+playerStatsForm.Strength,
		playerStats.Charisma+playerStatsForm.Charisma,
		playerStats.Constitution+playerStatsForm.Constitution,
		playerStats.Dexterity+playerStatsForm.Dexterity,
		playerStats.HP+playerStatsForm.HP,
		playerStats.Intelligence+playerStatsForm.Intelligence,
		playerStats.Mana+playerStatsForm.Mana,
		playerStats.Stamina+playerStatsForm.Stamina,
		playerStats.Wisdom+playerStatsForm.Wisdom)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	var newCharactersStats entities.Stats
	err = pgxscan.Get(ctx, db, &newCharactersStats, `SELECT character_id, strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma, hp FROM stats where character_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &newCharactersStats)

}
func UpdateCharacters(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	type playerInfoUpdatable struct {
		Exp           int
		Level         int
		InventorySize int
		Po            int
	}

	var playerInfoUpdatableForm playerInfoUpdatable
	if err := c.ShouldBindBodyWithJSON(&playerInfoUpdatableForm); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	selectedCharacters, err := repository.GetCharactersByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	_, err = db.Exec(ctx, `UPDATE player set (exp, level, inventory_size, po) = ($2,$3,$4,$5) where id = $1`, id, playerInfoUpdatableForm.Exp+selectedCharacters.Exp, playerInfoUpdatableForm.Level+selectedCharacters.Level, playerInfoUpdatableForm.InventorySize+selectedCharacters.InventorySize, playerInfoUpdatableForm.Po+selectedCharacters.Po)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
}
func UpdateCharactersLocation(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var locationID int
	if err := c.ShouldBindBodyWithJSON(&locationID); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	// check if location exist
	var locations entities.Locations
	err := pgxscan.Get(ctx, db, &locations, `SELECT id FROM locations where id = $1`, locationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	player, err := repository.GetCharactersByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// move player
	_, err = db.Exec(ctx, `UPDATE players SET location_id = $2 where id = $1`, player.ID, locationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	err = pgxscan.Get(ctx, db, &player, `SELECT id, email, username, race_id, job_id, exp, level, guild_id, inventory_size, po, location_id, user_id FROM characters where id = $1`, player.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &player)

}
func UpdateCharactersEquipment(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	type body struct {
		ItemID      int
		Emplacement string
	}

	var playerEquipmentForm body
	if err := c.ShouldBindBodyWithJSON(&playerEquipmentForm); err != nil {
		c.String(http.StatusBadRequest, "bad request during binding body")
		return
	}

	player, err := repository.GetCharactersByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var selectedItem entities.Item
	err = pgxscan.Get(ctx, db, &selectedItem, `SELECT id FROM items where id = $1`, playerEquipmentForm.ItemID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request selecting item")
		return
	}

	var playerEquipments entities.Equipment
	err = pgxscan.Get(ctx, db, &playerEquipments, `SELECT character_id, helmet, chestplate, leggings, boots, mainhand, offhand, accessory_slot_0, accessory_slot_1, accessory_slot_2, accessory_slot_3 FROM equipment where character_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request during getting current equipment")
		return
	}

	switch playerEquipmentForm.Emplacement {
	case "Helmet":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Helmet") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.Helmet, playerEquipments.CharactersId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.CharactersId, "helmet", db, c)
	case "Chestplate":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Chestplate") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.Chestplate, playerEquipments.CharactersId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.CharactersId, "chestplate", db, c)
	case "Leggings":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Leggings") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.Leggings, playerEquipments.CharactersId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.CharactersId, "leggings", db, c)
	case "Boots":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Boots") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.Boots, playerEquipments.CharactersId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.CharactersId, "boots", db, c)
	case "Mainhand":
		if repository.CheckEquipmentEmplacement(playerEquipments, "MainHand") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.Mainhand, playerEquipments.CharactersId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.CharactersId, "mainhand", db, c)
	case "Offhand":
		if repository.CheckEquipmentEmplacement(playerEquipments, "OffHand") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.Offhand, playerEquipments.CharactersId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.CharactersId, "offhand", db, c)
	case "AccessorySlot0":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Accesory0") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.AccessorySlot0, playerEquipments.CharactersId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.CharactersId, "accessory_slot_0", db, c)
	case "AccessorySlot1":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Accesory1") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.AccessorySlot1, playerEquipments.CharactersId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.CharactersId, "accessory_slot_1", db, c)
	case "AccessorySlot2":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Accesory2") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.AccessorySlot2, playerEquipments.CharactersId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.CharactersId, "accessory_slot_2", db, c)
	case "AccessorySlot3":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Accesory3") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.AccessorySlot3, playerEquipments.CharactersId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.CharactersId, "accessory_slot_3", db, c)
	default:
		break
	}

	var newCharactersEquipment entities.Equipment
	err = pgxscan.Get(ctx, db, &newCharactersEquipment, `SELECT character_id, helmet, chestplate, leggings, boots, mainhand, offhand, accessory_slot_0, accessory_slot_1, accessory_slot_2, accessory_slot_3 FROM equipment where character_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request during getting current equipment")
		return
	}

	c.JSON(http.StatusOK, &newCharactersEquipment)
}
func UpdateCharactersInventory(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	type body struct {
		itemID   int
		quantity int
	}

	var playerItemForm body
	if err := c.ShouldBindBodyWithJSON(&playerItemForm); err != nil {
		c.String(http.StatusBadRequest, "bad request during binding body")
		return
	}

	player, err := repository.GetCharactersByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var selectedItem entities.Item
	err = pgxscan.Get(ctx, db, &selectedItem, `SELECT id FROM items where id = $1`, playerItemForm.itemID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request selecting item")
		return
	}

	repository.DoUpsertItemInInventory(ctx, playerItemForm.itemID, player.ID, playerItemForm.quantity, db, c)

	var playerInventory []entities.Inventory
	err = pgxscan.Select(ctx, db, &playerInventory, `SELECT item_id, quantity FROM inventory where character_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerInventory)
}
func UpdateCharactersSkills(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	type Body struct {
		SkillID int
	}

	var playerSkillForm Body
	err := c.ShouldBindBodyWithJSON(&playerSkillForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// check if player exist
	player, err := repository.GetCharactersByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	// check if skill exist
	var selectedSkill entities.Skill
	err = pgxscan.Get(ctx, db, &selectedSkill, `SELECT id, name FROM skills where id = $1`, playerSkillForm.SkillID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request selecting pets")
		return
	}
	// check if player already have this pet
	var playerSkills entities.CharactersSkill
	err = pgxscan.Get(ctx, db, &playerSkills, `SELECT skill_id FROM character_skill where skill_id = $1`, playerSkillForm.SkillID)
	if errors.Is(err, pgx.ErrNoRows) {
		_, insertErr := db.Exec(ctx, `INSERT INTO character_skill (skill_id, character_id) values ($1,$2)`, playerSkillForm.SkillID, player.ID)
		if insertErr != nil {
			c.String(http.StatusBadRequest, "bad request inserting user's skill")
			return
		}
		c.JSON(http.StatusOK, "New SKill Added")
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		c.String(http.StatusBadRequest, "bad request selecting user's skill")
		return
	}

	if playerSkills.SkillID != 0 {

		var skillScroll entities.Item
		err = pgxscan.Get(ctx, db, &skillScroll, `SELECT id FROM items where name like '%$1%'`, selectedSkill.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}

		repository.DoUpsertItemInInventory(ctx, skillScroll.ID, player.ID, 1, db, c)

		c.JSON(http.StatusOK, "You Already have this skill, you got his scroll")
	}

}

// DELETE \\

func DeleteCharacters(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	// check if player exist
	player, err := repository.GetCharactersByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request during getting Characters"+err.Error())
		return
	}

	// delete pets
	_, err = db.Exec(ctx, `DELETE from character_pets_mounts  where character_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request deleting player's pets")
		return
	}
	// delete skills
	_, err = db.Exec(ctx, `DELETE from character_skill where character_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad  request deleting player's skills")
		return
	}
	// delete inventory
	_, err = db.Exec(ctx, `DELETE from inventory where character_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad  request deleting player's inventory")
		return
	}
	// delete equipment
	_, err = db.Exec(ctx, `DELETE from equipment where character_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad  request deleting player's equipment")
		return
	}
	// delete stat
	_, err = db.Exec(ctx, `DELETE from stats where character_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad  request deleting player's stats")
		return
	}
	// delete player action

	_, err = db.Exec(ctx, `DELETE from character_actions where character_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad  request deleting player's action")
		return
	}

	// delete player
	_, err = db.Exec(ctx, `DELETE from characters where id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad  request deleting player")
		return
	}
}
func DeleteCharactersItemInInventory(c *gin.Context) {
	db := database.Connect()
	playerID := c.Param("player")
	itemID := c.Param("item")

	item, err := repository.GetItemByID(ctx, db, itemID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	player, err := repository.GetCharactersByID(ctx, db, playerID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	_, deleteErr := db.Exec(ctx, `DELETE FROM inventory where item_id = $1 and character_id = $2`, item.ID, player.ID)
	if deleteErr != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.Status(http.StatusOK)

}
func DeleteCharactersPets(c *gin.Context) {
	db := database.Connect()
	playerID := c.Param("player")
	petID := c.Param("pet")

	pet, err := repository.GetPetByID(ctx, db, petID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	player, err := repository.GetCharactersByID(ctx, db, playerID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	_, deleteErr := db.Exec(ctx, `DELETE FROM character_pets_mounts where pet_id = $1 and character_id = $2`, pet.ID, player.ID)
	if deleteErr != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.Status(http.StatusOK)
}
func DeleteCharactersSkill(c *gin.Context) {
	db := database.Connect()
	playerID := c.Param("player")
	skillID := c.Param("skill")

	skill, err := repository.GetSkillInfoByID(ctx, db, skillID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	player, err := repository.GetCharactersByID(ctx, db, playerID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	_, deleteErr := db.Exec(ctx, `DELETE FROM character_skill where skill_id = $1 and character_id = $2`, skill.ID, player.ID)
	if deleteErr != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.Status(http.StatusOK)
}
