package controllers

import (
	"errors"
	"github.com/jackc/pgx/v5"
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

	creature, err := repository.GetCreatureByID(ctx, db, id)
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

	creature, err := repository.GetCreatureByID(ctx, db, id)
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

	creature, err := repository.GetCreatureByID(ctx, db, id)
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
	db := database.Connect()
	var CreatureForm entities.Creatures
	if err := c.ShouldBindBodyWithJSON(&CreatureForm); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	// check if creature exist
	existingCreature, ecErr := repository.GetCreatureByName(ctx, db, CreatureForm.Name)
	if errors.Is(ecErr, pgx.ErrNoRows) {
		c.JSON(http.StatusConflict, &existingCreature)
	}
	if ecErr != nil && !errors.Is(ecErr, pgx.ErrNoRows) {
		c.String(http.StatusBadRequest, "bad request")
	}
	// insert creature
	_, err := db.Exec(ctx, `INSERT INTO creatures (name, is_pet, strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma, level, hp) values ($1, $2, $3, $4,$5, $6, $7, $8,$9, $10, $11, $12)`,
		CreatureForm.Name, CreatureForm.IsPet, CreatureForm.Strength, CreatureForm.Constitution, CreatureForm.Mana, CreatureForm.Stamina,
		CreatureForm.Dexterity, CreatureForm.Intelligence, CreatureForm.Wisdom, CreatureForm.Charisma, CreatureForm.Level, CreatureForm.HP)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	// Get Creature
	creature, cErr := repository.GetCreatureByName(ctx, db, CreatureForm.Name)
	if cErr != nil {
		c.String(http.StatusBadRequest, "bad request")
	}
	// return Creature
	c.JSON(http.StatusCreated, &creature)
	c.Done()

}
func AddCreatureSpawn(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	var creatureSpawnForm entities.CreatureSpawns
	if err := c.ShouldBindBodyWithJSON(&creatureSpawnForm); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	creature, err := repository.GetCreatureByID(ctx, db, id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	var seletecdEmplacement entities.Locations
	err = pgxscan.Get(ctx, db, &seletecdEmplacement, `SELECT id FROM locations where id = $1`, creatureSpawnForm.EmplacementID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	_, err = db.Exec(ctx, `INSERT INTO creaturespawn (creature_id, emplacement_id, level_required, spawn_rate) values ($1, $2, $3, $4)`,
		creature.ID, seletecdEmplacement.ID, creatureSpawnForm.LevelRequired, creatureSpawnForm.SpawnRate)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	var creatureSpawn entities.CreatureSpawns
	err = pgxscan.Get(ctx, db, &creatureSpawn, `SELECT * FROM creaturespawn where creature_id = $1 AND emplacement_id = $2`, creature.ID, seletecdEmplacement.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusCreated, &creatureSpawn)
	c.Done()

}
func AddCreatureSkill(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	var creatureSkillForm entities.CreatureSkill
	if err := c.ShouldBindBodyWithJSON(&creatureSkillForm); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	creature, err := repository.GetCreatureByID(ctx, db, id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	var seletecdSkill entities.Locations
	err = pgxscan.Get(ctx, db, &seletecdSkill, `SELECT id FROM skills where id = $1`, creatureSkillForm.SkillID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	_, err = db.Exec(ctx, `INSERT INTO creature_skill (creature_id, skill_id) values ($1,$2)`, id, creatureSkillForm.SkillID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	var creatureSkill entities.CreatureSpawns
	err = pgxscan.Get(ctx, db, &creatureSkill, `SELECT * FROM skills s join creature_skill cs on s.id = cs.skill_id where creature_id = $1 and s.id = $2`, creature.ID, seletecdSkill.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusCreated, &creatureSkill)
	c.Done()
}
func AddCreatureLoot(c *gin.Context) {

}

// PATCH \\

func UpdateCreature(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var creatureUpdateForm entities.Creatures
	if err := c.ShouldBindBodyWithJSON(&creatureUpdateForm); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	creature, err := repository.GetCreatureByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	_, err = db.Exec(ctx, `UPDATE creatures set (name, is_pet, strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma, level, hp ) 
    = ($2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) where id = $1`,
		creature.ID, creatureUpdateForm.Name, creatureUpdateForm.IsPet, creatureUpdateForm.Strength, creatureUpdateForm.Mana, creatureUpdateForm.Stamina,
		creatureUpdateForm.Dexterity, creatureUpdateForm.Intelligence, creatureUpdateForm.Wisdom, creatureUpdateForm.Level, creatureUpdateForm.HP)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	creature, err = repository.GetCreatureByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &creature)
	c.Done()

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
