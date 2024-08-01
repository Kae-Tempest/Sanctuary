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
func GetAllPlayers(c *gin.Context) {
	db := database.Connect()
	var players []entities.Player
	err := pgxscan.Select(ctx, db, &players, `SELECT * FROM players`)
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
func GetOnePlayer(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	player, err := repository.GetPlayerByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &player)
}
func GetPlayerStats(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	player, err := repository.GetPlayerByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var playerStats entities.Stats
	err = pgxscan.Get(ctx, db, &playerStats, `SELECT * from stats where player_id = $1`, player.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerStats)

}
func GetPlayerEquipment(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	player, err := repository.GetPlayerByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var playerEquipment entities.Equipment
	err = pgxscan.Get(ctx, db, &playerEquipment, `SELECT * from equipment where player_id = $1`, player.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerEquipment)

}
func GetPlayerInventory(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	player, err := repository.GetPlayerByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var playerInventory entities.Inventory
	err = pgxscan.Get(ctx, db, &playerInventory, `SELECT * from inventory where player_id = $1`, player.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerInventory)

}
func GetPlayerPets(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	player, err := repository.GetPlayerByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var playerPets []entities.PlayerPet
	err = pgxscan.Select(ctx, db, &playerPets, `SELECT pet_id FROM player_pets_mounts where player_id = $1`, player.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerPets)

}
func GetPlayerGuild(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	player, err := repository.GetPlayerByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var playerGuild entities.Guild
	err = pgxscan.Get(ctx, db, &playerGuild, `SELECT * FROM guilds g join guilds_members gm on g.id = gm.guilds_id where gm.player_id = $1`, player.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerGuild)

}
func GetPlayerSkill(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	player, err := repository.GetPlayerByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var playerSkill []entities.Skill
	err = pgxscan.Select(ctx, db, &playerSkill, `SELECT * from player_skill where player_id = $1`, player.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerSkill)

}

// POST \\
func CreatePlayer(c *gin.Context) {
	db := database.Connect()
	var playerForm entities.Player
	if err := c.ShouldBindBodyWithJSON(&playerForm); err != nil {
		c.String(http.StatusBadRequest, "bad request during bind JSON")
		return
	}

	existingPlayer, epErr := repository.GetPlayerByEmail(ctx, db, playerForm.Email)
	if errors.Is(pgx.ErrNoRows, epErr) {
		c.JSON(http.StatusConflict, "existing player"+epErr.Error())
		fmt.Println(existingPlayer)
	}
	if epErr != nil && !errors.Is(epErr, pgx.ErrNoRows) {
		c.String(http.StatusBadRequest, "bad request")
	}

	_, err := db.Exec(ctx, `INSERT INTO players (email, username, race_id, job_id, exp, level, guild_id, inventory_size, po, location_id, user_id) values ($1, $2, $3, $4, 0, 1, 0, 10, 50, 1, null)`,
		playerForm.Email, playerForm.Username, playerForm.RaceID, playerForm.JobID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request during inserting data"+err.Error())
		return
	}
	var playerID int
	err = pgxscan.Get(ctx, db, &playerID, `SELECT id from players where email = $1`, playerForm.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request during get PlayerID"+err.Error())
		return
	}
	// Get Job Stats
	var playerJob entities.Job
	err = pgxscan.Get(ctx, db, &playerJob, `SELECT * from jobs where id = $1`, playerForm.JobID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request during get Job"+err.Error())
		return
	}
	// Get Race Stats
	var playerRace entities.Race
	err = pgxscan.Get(ctx, db, &playerRace, `SELECT * from races where id = $1`, playerForm.RaceID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request during get Race"+err.Error())
		return
	}
	// Create Player Stats
	_, err = db.Exec(ctx, `INSERT INTO stats (player_id, strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma, hp) values ($1, $2, $3, $4,$5,$6,$7,$8,$9,20)`,
		playerID, playerJob.Strength, playerJob.Constitution, playerRace.Mana, playerRace.Stamina, playerJob.Dexterity, playerJob.Intelligence, playerRace.Wisdom, playerRace.Charisma)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request during insert Stat"+err.Error())
		return
	}

	_, err = db.Exec(ctx, `INSERT INTO equipment (player_id, helmet, chestplate, leggings, boots, mainhand, offhand, accessory_slot_0, accessory_slot_1, accessory_slot_2, accessory_slot_3) VALUES 
    ($1,0,0,0,0,0,0,0,0,0,0)`, playerID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request during insert equipment"+err.Error())
		return
	}

	c.Status(http.StatusCreated)
}
func AddItemToPlayerInventory(c *gin.Context) {
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
	err := pgxscan.Get(ctx, db, &selectedItem, `SELECT * FROM items where id = $1`, body.ItemID)
	if err != nil {
		c.JSON(http.StatusNotFound, "Item with this ID doesn't exist !")
		return
	}

	player, err := repository.GetPlayerByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// verif if user had already item
	var playerInventory []entities.Inventory
	err = pgxscan.Select(ctx, db, &playerInventory, `SELECT item_id, quantity FROM inventory where player_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	if len(playerInventory) <= 0 {
		_, err = db.Exec(ctx, `INSERT into inventory (player_id, item_id, quantity) values ($1,$2,$3)`, player.ID, body.ItemID, body.Quantity)
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
		_, err = db.Exec(ctx, `INSERT into inventory (player_id, item_id, quantity) values ($1,$2,$3)`, player.ID, body.ItemID, body.Quantity)
		if err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		c.JSON(http.StatusCreated, "Item add to inventory")
	}
}
func AddPetToPlayer(c *gin.Context) {
	type Body struct {
		PetId int
	}
	db := database.Connect()
	id := c.Param("id")
	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var selectedPet entities.PetsMounts
	err := pgxscan.Get(ctx, db, &selectedPet, `SELECT id FROM pets_mounts where id = $1`, body.PetId)
	if err != nil {
		c.JSON(http.StatusNotFound, "Pet with this ID doesn't exist !")
		return
	}

	player, err := repository.GetPlayerByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var playerPets []entities.PlayerPet
	err = pgxscan.Select(ctx, db, &playerPets, `SELECT pet_id FROM player_pets_mounts where player_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	if len(playerPets) <= 0 {
		_, err = db.Exec(ctx, `INSERT into player_pets_mounts (player_id, pet_id) values ($1,$2)`, player.ID, body.PetId)
		if err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		c.JSON(http.StatusCreated, "First pet get by User")
	}

	var havePet = false
	for _, pet := range playerPets {
		if pet.PetID == body.PetId {
			havePet = true
			break
		} else {
			havePet = false
			continue
		}
	}

	if havePet {
		c.JSON(http.StatusBadRequest, "User already have this pet")
	} else {
		_, err = db.Exec(ctx, `INSERT into player_pets_mounts (player_id, pet_id) values ($1,$2)`, player.ID, body.PetId)
		if err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		c.JSON(http.StatusCreated, "User tamed this pet")
	}

}
func AddSkillToPlayer(c *gin.Context) {
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

	player, err := repository.GetPlayerByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var playerSkills []entities.PlayerSkill
	err = pgxscan.Select(ctx, db, &playerSkills, `SELECT skill_id FROM player_skill where player_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	if len(playerSkills) <= 0 {
		_, err = db.Exec(ctx, `INSERT into player_skill (player_id, skill_id) values ($1,$2)`, player.ID, body.SkillId)
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
		_, err = db.Exec(ctx, `INSERT into player_skill (player_id, skill_id) values ($1,$2)`, player.ID, body.SkillId)
		if err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		c.JSON(http.StatusCreated, "User tamed this skill")
	}

}

// PATCH \\
func UpdatePlayerStats(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	var playerStatsForm entities.Stats
	if err := c.ShouldBindBodyWithJSON(&playerStatsForm); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	player, err := repository.GetPlayerByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var playerStats entities.Stats
	err = pgxscan.Get(ctx, db, &playerStats, `SELECT * FROM stats where player_id = $1`, player.ID)
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

	var newPlayerStats entities.Stats
	err = pgxscan.Get(ctx, db, &newPlayerStats, `SELECT * FROM stats where player_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &newPlayerStats)

}
func UpdatePlayer(c *gin.Context) {
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

	selectedPlayer, err := repository.GetPlayerByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	_, err = db.Exec(ctx, `UPDATE player set (exp, level, inventory_size, po) = ($2,$3,$4,$5) where id = $1`, id, playerInfoUpdatableForm.Exp+selectedPlayer.Exp, playerInfoUpdatableForm.Level+selectedPlayer.Level, playerInfoUpdatableForm.InventorySize+selectedPlayer.InventorySize, playerInfoUpdatableForm.Po+selectedPlayer.Po)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
}
func UpdatePlayerLocation(c *gin.Context) {
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

	player, err := repository.GetPlayerByID(ctx, db, id)
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

	err = pgxscan.Get(ctx, db, &player, `SELECT * FROM players where id = $1`, player.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &player)

}
func UpdatePlayerEquipment(c *gin.Context) {
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

	player, err := repository.GetPlayerByID(ctx, db, id)
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
	err = pgxscan.Get(ctx, db, &playerEquipments, `SELECT * FROM equipment where player_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request during getting current equipment")
		return
	}

	switch playerEquipmentForm.Emplacement {
	case "Helmet":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Helmet") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.Helmet, playerEquipments.PlayerId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.PlayerId, "helmet", db, c)
	case "Chestplate":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Chestplate") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.Chestplate, playerEquipments.PlayerId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.PlayerId, "chestplate", db, c)
	case "Leggings":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Leggings") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.Leggings, playerEquipments.PlayerId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.PlayerId, "leggings", db, c)
	case "Boots":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Boots") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.Boots, playerEquipments.PlayerId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.PlayerId, "boots", db, c)
	case "Mainhand":
		if repository.CheckEquipmentEmplacement(playerEquipments, "MainHand") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.Mainhand, playerEquipments.PlayerId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.PlayerId, "mainhand", db, c)
	case "Offhand":
		if repository.CheckEquipmentEmplacement(playerEquipments, "OffHand") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.Offhand, playerEquipments.PlayerId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.PlayerId, "offhand", db, c)
	case "AccessorySlot0":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Accesory0") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.AccessorySlot0, playerEquipments.PlayerId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.PlayerId, "accessory_slot_0", db, c)
	case "AccessorySlot1":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Accesory1") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.AccessorySlot1, playerEquipments.PlayerId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.PlayerId, "accessory_slot_1", db, c)
	case "AccessorySlot2":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Accesory2") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.AccessorySlot2, playerEquipments.PlayerId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.PlayerId, "accessory_slot_2", db, c)
	case "AccessorySlot3":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Accesory3") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.AccessorySlot3, playerEquipments.PlayerId, 1, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.ItemID, playerEquipments.PlayerId, "accessory_slot_3", db, c)
	default:
		break
	}

	var newPlayerEquipment entities.Equipment
	err = pgxscan.Get(ctx, db, &newPlayerEquipment, `SELECT * FROM equipment where player_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request during getting current equipment")
		return
	}

	c.JSON(http.StatusOK, &newPlayerEquipment)
}
func UpdatePlayerInventory(c *gin.Context) {
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

	player, err := repository.GetPlayerByID(ctx, db, id)
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
	err = pgxscan.Select(ctx, db, &playerInventory, `SELECT item_id, quantity FROM inventory where player_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerInventory)
}
func UpdatePlayerPets(c *gin.Context) {
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
	player, err := repository.GetPlayerByID(ctx, db, id)
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
	var playerPets entities.PlayerPet
	err = pgxscan.Get(ctx, db, &playerPets, `SELECT pet_id FROM player_pets_mounts where pet_id = $1`, playerPetForm.PetsID)
	if errors.Is(err, pgx.ErrNoRows) {
		_, insertErr := db.Exec(ctx, `INSERT INTO player_pets_mounts (pet_id, player_id) values ($1,$2)`, playerPetForm.PetsID, player.ID)
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
		err = pgxscan.Get(ctx, db, &petScroll, `SELECT * FROM items where name like '%$1%'`, mob.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}

		repository.DoUpsertItemInInventory(ctx, petScroll.ID, player.ID, 1, db, c)

		c.JSON(http.StatusOK, "You Already have this pet, you got his scroll")
	}
}
func UpdatePlayerSkills(c *gin.Context) {
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
	player, err := repository.GetPlayerByID(ctx, db, id)
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
	var playerSkills entities.PlayerSkill
	err = pgxscan.Get(ctx, db, &playerSkills, `SELECT skill_id FROM player_skill where skill_id = $1`, playerSkillForm.SkillID)
	if errors.Is(err, pgx.ErrNoRows) {
		_, insertErr := db.Exec(ctx, `INSERT INTO player_skill (skill_id, player_id) values ($1,$2)`, playerSkillForm.SkillID, player.ID)
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
		err = pgxscan.Get(ctx, db, &skillScroll, `SELECT * FROM items where name like '%$1%'`, selectedSkill.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}

		repository.DoUpsertItemInInventory(ctx, skillScroll.ID, player.ID, 1, db, c)

		c.JSON(http.StatusOK, "You Already have this skill, you got his scroll")
	}

}

// DELETE \\

func DeletePlayer(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	// check if player exist
	player, err := repository.GetPlayerByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request during getting Player"+err.Error())
		return
	}

	// delete pets
	_, err = db.Exec(ctx, `DELETE from player_pets_mounts  where player_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request deleting player's pets")
		return
	}
	// delete skills
	_, err = db.Exec(ctx, `DELETE from player_skill where player_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad  request deleting player's skills")
		return
	}
	// delete inventory
	_, err = db.Exec(ctx, `DELETE from inventory where player_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad  request deleting player's inventory")
		return
	}
	// delete equipment
	_, err = db.Exec(ctx, `DELETE from equipment where player_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad  request deleting player's equipment")
		return
	}
	// delete stat
	_, err = db.Exec(ctx, `DELETE from stats where player_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad  request deleting player's stats")
		return
	}
	// delete player action

	_, err = db.Exec(ctx, `DELETE from players_actions where player_id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad  request deleting player's action")
		return
	}

	// delete player
	_, err = db.Exec(ctx, `DELETE from players where id = $1`, player.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad  request deleting player")
		return
	}
}
func DeletePlayerItemInInventory(c *gin.Context) {
	db := database.Connect()
	playerID := c.Param("player")
	itemID := c.Param("item")

	item, err := repository.GetItemByID(ctx, db, itemID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	player, err := repository.GetPlayerByID(ctx, db, playerID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	_, deleteErr := db.Exec(ctx, `DELETE FROM inventory where item_id = $1 and player_id = $2`, item.ID, player.ID)
	if deleteErr != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.Status(http.StatusOK)

}
func DeletePlayerPets(c *gin.Context) {
	db := database.Connect()
	playerID := c.Param("player")
	petID := c.Param("pet")

	pet, err := repository.GetPetByID(ctx, db, petID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	player, err := repository.GetPlayerByID(ctx, db, playerID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	_, deleteErr := db.Exec(ctx, `DELETE FROM player_pets_mounts where pet_id = $1 and player_id = $2`, pet.ID, player.ID)
	if deleteErr != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.Status(http.StatusOK)
}
func DeletePlayerSkill(c *gin.Context) {
	db := database.Connect()
	playerID := c.Param("player")
	skillID := c.Param("skill")

	skill, err := repository.GetSkillInfoByID(ctx, db, skillID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	player, err := repository.GetPlayerByID(ctx, db, playerID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	_, deleteErr := db.Exec(ctx, `DELETE FROM player_skill where skill_id = $1 and player_id = $2`, skill.ID, player.ID)
	if deleteErr != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.Status(http.StatusOK)
}
