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
)

var ctx = context.Background()

// GET \\
func GetAllPlayers(c *gin.Context) {
	db := database.Connect()
	var players []entities.Player
	err := pgxscan.Select(ctx, db, &players, `SELECT * FROM players`)
	if err != nil {
		slog.Error("Error during selection all player", err)
	}
	fmt.Println(&players)
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
	}
	c.JSON(http.StatusOK, &playerInventory)
	c.Done()
}
func GetPlayerPets(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	var playerPets []entities.UserPet
	err := pgxscan.Select(ctx, db, &playerPets, `SELECT pet_id from user_pets_mounts where player_id = $1`, id)
	if err != nil {
		slog.Error("Error during selection player pet with id:"+id, err)
	}
	c.JSON(http.StatusOK, &playerPets)
	c.Done()
}
func GetPlayerGuild(c *gin.Context) {

}
func GetPlayerSkill(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	var playerSkill []entities.Skill
	err := pgxscan.Select(ctx, db, &playerSkill, `SELECT * from user_skill where player_id = $1`, id)
	if err != nil {
		slog.Error("Error during selection player Skill with id:"+id, err)
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
	fmt.Println(playerForm)
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
	err = pgxscan.Select(ctx, db, &playerInventory, `SELECT * FROM inventory where player_id = $1`, playerID)
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

}
func AddSkillToPlayer(c *gin.Context) {

}
