package controllers

import (
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"sanctuary-api/database"
	"sanctuary-api/entities"
	"sanctuary-api/repository"
	"strconv"
)

func GetLocations(c *gin.Context) {
	db := database.Connect()

	var locations []entities.Locations
	err := pgxscan.Select(ctx, db, &locations, `SELECT id, name, is_safety, difficulty, type, size FROM locations`)
	if err != nil {
		slog.Error("Error during selecting locations", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during selecting locations")
		return
	}

	c.JSON(http.StatusOK, &locations)
}
func GetLocationByID(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var location entities.Locations
	err := pgxscan.Select(ctx, db, &location, `SELECT id, name, is_safety, difficulty, type, size FROM locations where id = $1`, id)
	if err != nil {
		slog.Error("Error during selecting location by ID", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during selecting location by ID")
		return
	}

	c.JSON(http.StatusOK, &location)

}
func GetCharactersByLocation(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	players, err := repository.GetCharactersByLocation(ctx, db, id)
	if err != nil {
		slog.Error("Error during selecting characters by is location", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during selecting characters by is location")
		return
	}

	c.JSON(http.StatusOK, &players)
}
func GetCreaturesByLocation(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	creatures, err := repository.GetMobsByLocation(ctx, db, id)
	if err != nil {
		slog.Error("Error during selecting mobs by is location", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during selecting mobs by is location")
		return
	}

	c.JSON(http.StatusOK, &creatures)
}
func GetResourcesByLocation(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var resources []entities.Resources
	err := pgxscan.Select(ctx, db, &resources, `SELECT id, name, location_id, item_id, quantities_per_min FROM resources where location_id = $1`, id)
	if err != nil {
		slog.Error("Error during resources mob by is location", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during selecting resources by is location")
		return
	}

	c.JSON(http.StatusOK, &resources)
}

func CreateLocation(c *gin.Context) {
	db := database.Connect()

	var locationForm entities.Locations
	if err := c.ShouldBindBodyWithJSON(&locationForm); err != nil {
		slog.Error("Error during binding locationForm", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during binding locationForm")
		return
	}

	location, err := repository.GetLocationByName(ctx, db, locationForm.Name)
	if err != nil {
		slog.Error("Error during selecting location by name", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during selecting location by name")
		return
	}

	if location.Name == locationForm.Name {
		slog.Error("Location Name already existing", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Location Name already existing")
		return
	} else {
		_, insertErr := db.Exec(ctx, `INSERT INTO locations (name, is_safety, difficulty, type, size) values ($1,$2,$3,$4,$5)`,
			locationForm.Name, locationForm.IsSafety, locationForm.Difficulty, locationForm.Type, locationForm.Size)

		if insertErr != nil {
			slog.Error("Error during inserting location", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, "Error during inserting location")
			return
		}

		newLocation, selectErr := repository.GetLocationByName(ctx, db, locationForm.Name)
		if selectErr != nil {
			slog.Error("Error during selecting location by name", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, "Error during selecting location by name")
			return
		}

		c.JSON(http.StatusCreated, &newLocation)

	}
}

func UpdateLocation(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var locationForm entities.Locations
	if err := c.ShouldBindBodyWithJSON(&locationForm); err != nil {
		slog.Error("Error during binding locationForm", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during binding locationForm")
		return
	}
	location, err := repository.GetLocationByID(ctx, db, id)
	if err != nil {
		slog.Error("Error during selecting location by ID", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during selecting location by ID")
		return
	}

	_, updateErr := db.Exec(ctx, `UPDATE locations set (name, is_safety, difficulty, type, size) = ($2,$3,$4,$5,$6) where id = $1`,
		location.ID, locationForm.IsSafety, locationForm.Difficulty, locationForm.Type, locationForm.Size)
	if updateErr != nil {
		slog.Error("Error during updating location", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during updating location")
		return
	}

	updateLocation, err := repository.GetLocationByID(ctx, db, strconv.Itoa(location.ID))
	if err != nil {
		slog.Error("Error during selecting location by ID", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during selecting location by ID")
		return
	}

	c.JSON(http.StatusOK, &updateLocation)
}

func DeleteLocation(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	location, err := repository.GetLocationByID(ctx, db, id)
	if err != nil {
		slog.Error("Error during selecting location by ID", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during selecting location by ID")
		return
	}

	players, err := repository.GetCharactersByLocation(ctx, db, id)
	if err != nil {
		slog.Error("Error during selecting character by is location", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during selecting character by is location")
		return
	}

	newLocation, locationErr := repository.GetLocationByName(ctx, db, "Lost Land")
	if locationErr != nil {
		slog.Error("Error during selecting location by Name", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during selecting location by Name")
		return
	}

	updateErr := repository.UpdateCharactersLocation(ctx, db, newLocation.ID, players)
	if updateErr != nil {
		slog.Error("Error during updating character's location", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during selecting character's by is location")
		return
	}

	creatures, cErr := repository.GetMobsByLocation(ctx, db, id)
	if cErr != nil {
		slog.Error("Error during selecting mobs by is location", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during selecting mobs by is location")
		return
	}
	updateErr = repository.UpdateMobsLocation(ctx, db, 0, creatures)

	_, deleteErr := db.Exec(ctx, `DELETE FROM locations where id = $1`, location.ID)
	if deleteErr != nil {
		slog.Error("Error during updating mob's location", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during selecting mob's by is location")
		return
	}

	c.Status(http.StatusOK)
}
