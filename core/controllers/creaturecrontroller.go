package controllers

import (
	"net/http"
	"sanctuary-api/database"
	"sanctuary-api/entities"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
)

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

	var creature entities.Creatures
	err := pgxscan.Get(ctx, db, &creature, `SELECT * FROM creatures where id = $1`, id)
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

	var creatureSpawn entities.CreatureSpawns
	err := pgxscan.Get(ctx, db, &creatureSpawn, `SELECT * FROM creaturespawn where creature_id = $1`, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &creatureSpawn)
	c.Done()
}
