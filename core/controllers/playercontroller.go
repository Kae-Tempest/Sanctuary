package controllers

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"sanctuary-api/database"
	"sanctuary-api/entities"
	"sanctuary-api/repository"
)

var ctx = context.Background()

// GET \\
func GetAllPlayers(c *gin.Context) {
	db := database.Connect()
	var players []entities.Player
	err := pgxscan.Select(ctx, db, &players, `SELECT * FROM players`)
	if err != nil {
		slog.Error("Error during selection all player", err)
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	if len(players) > 0 {
		c.JSON(http.StatusOK, &players)
	} else {
		c.JSON(http.StatusNotFound, gin.H{})
	}
	c.Done()
}
func GetOnePlayer(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	var players entities.Player
	err := pgxscan.Get(ctx, db, &players, `SELECT * FROM players where ID = $1`, id)
	if err != nil {
		slog.Error("Error during selection player id:"+id, err)
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &players)
	c.Done()
}
func GetPlayerStats(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	var playerStats entities.Stats
	err := pgxscan.Get(ctx, db, &playerStats, `SELECT * from stats where player_id = $1`, id)
	if err != nil {
		slog.Error("Error during selection player stats with id:"+id, err)
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerStats)
	c.Done()
}
func GetPlayerEquipment(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	var playerEquipment entities.Equipment
	err := pgxscan.Get(ctx, db, &playerEquipment, `SELECT * from equipment where player_id = $1`, id)
	if err != nil {
		slog.Error("Error during selection player equipments with id:"+id, err)
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerEquipment)
	c.Done()
}
func GetPlayerInventory(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	var playerInventory entities.Inventory
	err := pgxscan.Get(ctx, db, &playerInventory, `SELECT * from inventory where player_id = $1`, id)
	if err != nil {
		slog.Error("Error during selection player equipments with id:"+id, err)
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerInventory)
	c.Done()
}
func GetPlayerPets(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	var playerPets []entities.UserPet
	err := pgxscan.Select(ctx, db, &playerPets, `SELECT pet_id FROM user_pets_mounts where player_id = $1`, id)
	if err != nil {
		slog.Error("Error during selection player pet with id:"+id, err)
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerPets)
	c.Done()
}
func GetPlayerGuild(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	var playerGuild entities.Guild
	err := pgxscan.Get(ctx, db, &playerGuild, `SELECT * FROM guilds g join guilds_members gm on g.id = gm.guilds_id where gm.user_id = $1`, id)
	if err != nil {
		slog.Error("Error during selection player guild with id:"+id, err)
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerGuild)
	c.Done()
}
func GetPlayerSkill(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	var playerSkill []entities.Skill
	err := pgxscan.Select(ctx, db, &playerSkill, `SELECT * from user_skill where player_id = $1`, id)
	if err != nil {
		slog.Error("Error during selection player Skill with id:"+id, err)
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	c.JSON(http.StatusOK, &playerSkill)
	c.Done()
}

// POST \\
func CreatePlayer(c *gin.Context) {
	db := database.Connect()
	var playerForm entities.Player
	if err := c.ShouldBindBodyWithJSON(&playerForm); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	_, err := db.Exec(ctx, `INSERT INTO players (email, username, race_id, job_id, exp, level, guild_id, inventory_size, po, location_id) values ($1, $2, $3, $4, 0, 1, 0, 10, 50, 1)`,
		playerForm.Email, playerForm.Username, playerForm.RaceID, playerForm.JobID)
	if err != nil {
		slog.Error("Error during inserting player into database", err)
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	var playerID int
	err = pgxscan.Select(ctx, db, playerID, `SELECT id from players where email = $1`, playerForm.Email)
	// Get Job Stats
	var playerJob entities.Job
	err = pgxscan.Get(ctx, db, &playerJob, `SELECT * from jobs where id = $1`, playerForm.JobID)
	if err != nil {
		slog.Error(fmt.Sprintf("Error during selection of jobs with id: %d", playerForm.JobID), err)
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	// Get Race Stats
	var playerRace entities.Race
	err = pgxscan.Get(ctx, db, &playerRace, `SELECT * from races where id = $1`, playerForm.RaceID)
	if err != nil {
		slog.Error(fmt.Sprintf("Error during selection of race with id: %d", playerForm.RaceID), err)
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	// Create Player Stats
	_, err = db.Exec(ctx, `INSERT INTO stats (player_id, strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma, hp) values ($1, $2, $3, $4,$5,$6,$7,$8,$9,20)`,
		playerID, playerJob.Strength, playerJob.Constitution, playerRace.Mana, playerRace.Stamina, playerJob.Dexterity, playerJob.Intelligence, playerRace.Wisdom, playerRace.Charisma)
	if err != nil {
		slog.Error("Error during inserting player stats into database", err)
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.Status(http.StatusCreated)
	c.Done()
}
func AddItemToPlayerInventory(c *gin.Context) {
	type Body struct {
		ItemID   int
		Quantity int
	}
	db := database.Connect()
	playerID := c.Param("id")
	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// verifier si l'item existe
	var selectedItem entities.Items
	err := pgxscan.Get(ctx, db, &selectedItem, `SELECT * FROM items where id = $1`, body.ItemID)
	if err != nil {
		slog.Error("Item doesn't  exist", err)
		c.JSON(http.StatusNotFound, "Item with this ID doesn't exist !")
		return
	}

	// verif if user had already item
	var playerInventory []entities.Inventory
	err = pgxscan.Select(ctx, db, &playerInventory, `SELECT item_id, quantity FROM inventory where player_id = $1`, playerID)
	if err != nil {
		slog.Error("Error during selection player inventory with id:"+playerID, err)
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	if len(playerInventory) <= 0 {
		_, err = db.Exec(ctx, `INSERT into inventory (player_id, item_id, quantity) values ($1,$2,$3)`, playerID, body.ItemID, body.Quantity)
		if err != nil {
			slog.Error("Error during insertion in player inventory with id:"+playerID, err)
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
		_, err = db.Exec(ctx, `UPDATE inventory SET quantity = $1 where player_id = $2`, quantity+body.Quantity, playerID)
		if err != nil {
			slog.Error("Error during updating in player inventory with id:"+playerID, err)
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		c.JSON(http.StatusCreated, "Item quantity updated")
	} else {
		_, err = db.Exec(ctx, `INSERT into inventory (player_id, item_id, quantity) values ($1,$2,$3)`, playerID, body.ItemID, body.Quantity)
		if err != nil {
			slog.Error("Error during insertion in player inventory with id:"+playerID, err)
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
	playerID := c.Param("id")
	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var selectedPet entities.PetsMounts
	err := pgxscan.Get(ctx, db, &selectedPet, `SELECT id FROM pets_mounts where id = $1`, body.PetId)
	if err != nil {
		slog.Error("Pet doesn't exist", err)
		c.JSON(http.StatusNotFound, "Pet with this ID doesn't exist !")
		return
	}

	var playerPets []entities.UserPet
	err = pgxscan.Select(ctx, db, &playerPets, `SELECT pet_id FROM user_pets_mounts where player_id = $1`, playerID)
	if err != nil {
		slog.Error("Error during selection player pets with id:"+playerID, err)
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	if len(playerPets) <= 0 {
		_, err = db.Exec(ctx, `INSERT into user_pets_mounts (player_id, pet_id) values ($1,$2)`, playerID, body.PetId)
		if err != nil {
			slog.Error("Error during insertion in player pets with id:"+playerID, err)
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
		_, err = db.Exec(ctx, `INSERT into user_pets_mounts (player_id, pet_id) values ($1,$2)`, playerID, body.PetId)
		if err != nil {
			slog.Error("Error during insertion in player pets with id:"+playerID, err)
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
	playerID := c.Param("id")
	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var selectedSkill entities.Skill
	err := pgxscan.Get(ctx, db, &selectedSkill, `SELECT id FROM skills where id = $1`, body.SkillId)
	if err != nil {
		slog.Error("Pet doesn't exist", err)
		c.JSON(http.StatusNotFound, "Pet with this ID doesn't exist !")
		return
	}

	var playerSkills []entities.UserSkill
	err = pgxscan.Select(ctx, db, &playerSkills, `SELECT skill_id FROM user_skill where player_id = $1`, playerID)
	if err != nil {
		slog.Error("Error during selection player pets with id:"+playerID, err)
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	if len(playerSkills) <= 0 {
		_, err = db.Exec(ctx, `INSERT into user_skill (player_id, skill_id) values ($1,$2)`, playerID, body.SkillId)
		if err != nil {
			slog.Error("Error during insertion in player skill with id:"+playerID, err)
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
		_, err = db.Exec(ctx, `INSERT into user_skill (player_id, skill_id) values ($1,$2)`, playerID, body.SkillId)
		if err != nil {
			slog.Error("Error during insertion in player skill with id:"+playerID, err)
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
	var player entities.Player
	err := pgxscan.Get(ctx, db, &player, `SELECT id FROM players where id = $1`, id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	var playerStats entities.Stats
	err = pgxscan.Get(ctx, db, &playerStats, `SELECT * FROM stats where player_id = $1`, id)
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
	err = pgxscan.Get(ctx, db, &newPlayerStats, `SELECT * FROM stats where player_id = $1`, id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &newPlayerStats)
	c.Done()

}
func UpdatePlayerEquipment(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	type body struct {
		itemID      int
		Emplacement string
	}

	var playerEquipmentForm body
	if err := c.ShouldBindBodyWithJSON(&playerEquipmentForm); err != nil {
		c.String(http.StatusBadRequest, "bad request during binding body")
		return
	}

	var player entities.Player
	err := pgxscan.Get(ctx, db, &player, `SELECT id FROM players where id = $1`, id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request selecting player")
		return
	}

	var selectedItem entities.Items
	err = pgxscan.Get(ctx, db, &selectedItem, `SELECT id FROM items where id = $1`, playerEquipmentForm.itemID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request selecting item")
		return
	}

	var playerEquipments entities.Equipment
	err = pgxscan.Get(ctx, db, &playerEquipments, `SELECT * FROM equipment where player_id = $1`, id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request during getting current equipment")
		return
	}

	switch playerEquipmentForm.Emplacement {
	case "Helmet":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Helmet") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.Helmet, playerEquipments.PlayerId, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.itemID, playerEquipments.PlayerId, "helmet", db, c)
	case "Chestplate":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Chestplate") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.Chestplate, playerEquipments.PlayerId, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.itemID, playerEquipments.PlayerId, "chestplate", db, c)
	case "Leggings":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Leggings") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.Leggings, playerEquipments.PlayerId, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.itemID, playerEquipments.PlayerId, "leggings", db, c)
	case "Boots":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Boots") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.Boots, playerEquipments.PlayerId, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.itemID, playerEquipments.PlayerId, "boots", db, c)
	case "Mainhand":
		if repository.CheckEquipmentEmplacement(playerEquipments, "MainHand") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.Mainhand, playerEquipments.PlayerId, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.itemID, playerEquipments.PlayerId, "mainhand", db, c)
	case "Offhand":
		if repository.CheckEquipmentEmplacement(playerEquipments, "OffHand") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.Offhand, playerEquipments.PlayerId, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.itemID, playerEquipments.PlayerId, "offhand", db, c)
	case "AccessorySlot0":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Accesory0") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.AccessorySlot0, playerEquipments.PlayerId, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.itemID, playerEquipments.PlayerId, "accessory_slot_0", db, c)
	case "AccessorySlot1":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Accesory1") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.AccessorySlot1, playerEquipments.PlayerId, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.itemID, playerEquipments.PlayerId, "accessory_slot_2", db, c)
	case "AccessorySlot2":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Accesory2") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.AccessorySlot2, playerEquipments.PlayerId, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.itemID, playerEquipments.PlayerId, "accessory_slot_3", db, c)
	case "AccessorySlot3":
		if repository.CheckEquipmentEmplacement(playerEquipments, "Accesory3") {
			repository.DoUpsertItemInInventory(ctx, playerEquipments.AccessorySlot3, playerEquipments.PlayerId, db, c)
		}
		repository.DoUpdateEquipment(ctx, playerEquipmentForm.itemID, playerEquipments.PlayerId, "accessory_slot_4", db, c)
	default:
		break
	}
}

func UpdatePlayerInventory(c *gin.Context) {}
func UpdatePlayerPets(c *gin.Context)      {}
func UpdatePlayerSkills(c *gin.Context)    {}

// DELETE \\

func DeletePlayer(c *gin.Context) {}
