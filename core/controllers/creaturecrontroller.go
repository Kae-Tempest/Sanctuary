package controllers

import (
	"net/http"
	"sanctuary-api/database"
	"sanctuary-api/entities"
	"sanctuary-api/repository"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
)

// GET \\

func GetAllCreatures(c *gin.Context) {
	db := database.Connect()
	var creatures []entities.Creatures
	err := pgxscan.Select(ctx, db, &creatures, `SELECT * FROM creatures`)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	if len(creatures) > 0 {
		c.JSON(http.StatusOK, &creatures)
	} else {
		c.JSON(http.StatusNotFound, gin.H{})
	}
	c.Done()
}

func GetOneCreature(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	creature, err := repository.GetCreatureById(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &creature)
	c.Done()
}

func GetCreatureSpawn(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	creature, err := repository.GetCreatureById(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var creatureSpawn entities.CreatureSpawns
	err = pgxscan.Get(ctx, db, &creatureSpawn, `SELECT * FROM creaturespawn where creature_id = $1`, creature.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &creatureSpawn)
	c.Done()
}

func GetCreatureSkill(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	creature, err := repository.GetCreatureById(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var creatureSkill []entities.CreatureSkill
	err = pgxscan.Select(ctx, db, &creatureSkill, `SELECT * FROM creature_skill where creature_id = $1`, creature.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &creatureSkill)
	c.Done()
}

// POST \\

func CreateCreature(c *gin.Context) {

}

func AddCreatureSpawn(c *gin.Context) {

}

func AddCreatureSkill(c *gin.Context) {

}

func AddCreatureLoot(c *gin.Context) {

}

// PATCH \\

func UpdateCreature(c *gin.Context) {

}

func UpdateCreatureSpawn(c *gin.Context) {

}

func UpdateCreatureSkill(c *gin.Context) {

}

func UpdateCreatureLoot(c *gin.Context) {

}

// DELETE \\

func DeleteCreature(c *gin.Context) {

}

func DeleteCreatureSpawn(c *gin.Context) {

}

func DeleteCreatureSkill(c *gin.Context) {

}

func DeleteCreatureLoot(c *gin.Context) {

}
