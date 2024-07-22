package controllers

import (
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"net/http"
	"sanctuary-api/database"
	"sanctuary-api/entities"
	"sanctuary-api/repository"
	"strconv"
)

func GetLocations(c *gin.Context) {
	db := database.Connect()

	var locations []entities.Locations
	err := pgxscan.Select(ctx, db, &locations, `SELECT * FROM locations`)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &locations)
	c.Done()
}
func GetLocationByID(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var location entities.Locations
	err := pgxscan.Select(ctx, db, &location, `SELECT * FROM locations where id = $1`, id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &location)
	c.Done()
}
func GetPlayersByLocation(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	players, err := repository.GetPlayersByLocation(ctx, db, id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &players)
	c.Done()
}
func GetCreaturesByLocation(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	creatures, err := repository.GetCreaturesByLocation(ctx, db, id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &creatures)
	c.Done()
}
func GetResourcesByLocation(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var resources []entities.Resources
	err := pgxscan.Select(ctx, db, &resources, `SELECT * FROM resources where emplacement_id = $1`, id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &resources)
	c.Done()
}

func CreateLocation(c *gin.Context) {
	db := database.Connect()

	var locationForm entities.Locations
	if err := c.ShouldBindBodyWithJSON(&locationForm); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	location, err := repository.GetLocationByName(ctx, db, locationForm.Name)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	if location.Name == locationForm.Name {
		c.String(http.StatusBadRequest, "bad request")
		return
	} else {
		_, insertErr := db.Exec(ctx, `INSERT INTO locations (name, is_safety, difficulty, type, size) values ($1,$2,$3,$4,$5)`,
			locationForm.Name, locationForm.IsSafety, locationForm.Difficulty, locationForm.Type, locationForm.Size)

		if insertErr != nil {
			c.String(http.StatusBadRequest, "bad request")
			return
		}

		newLocation, selectErr := repository.GetLocationByName(ctx, db, locationForm.Name)
		if selectErr != nil {
			c.String(http.StatusBadRequest, "bad request")
			return
		}

		c.JSON(http.StatusCreated, &newLocation)
		c.Done()
	}
}

func UpdateLocation(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var locationForm entities.Locations
	if err := c.ShouldBindBodyWithJSON(&locationForm); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	location, err := repository.GetLocationByID(ctx, db, id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	_, updateErr := db.Exec(ctx, `UPDATE locations set (name, is_safety, difficulty, type, size) = ($2,$3,$4,$5,$6) where id = $1`,
		location.ID, locationForm.IsSafety, locationForm.Difficulty, locationForm.Type, locationForm.Size)
	if updateErr != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	updateLocation, err := repository.GetLocationByID(ctx, db, strconv.Itoa(location.ID))
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &updateLocation)
	c.Done()
}

func DeleteLocation(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	location, err := repository.GetLocationByID(ctx, db, id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	players, err := repository.GetPlayersByLocation(ctx, db, id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	newLocation, locationErr := repository.GetLocationByName(ctx, db, "Lost Land")
	if locationErr != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	updateErr := repository.UpdatePlayersLocation(ctx, db, newLocation.ID, players)
	if updateErr != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	creatures, cErr := repository.GetCreaturesByLocation(ctx, db, id)
	if cErr != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	updateErr = repository.UpdateCreaturesLocation(ctx, db, 0, creatures)

	_, deleteErr := db.Exec(ctx, `DELETE FROM locations where id = $1`, location.ID)
	if deleteErr != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.Status(http.StatusOK)
	c.Done()
}
