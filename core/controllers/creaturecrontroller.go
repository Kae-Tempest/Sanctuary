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
	db := database.Connect()
	id := c.Param("id")

	var creatureSpawnForm entities.CreatureSpawns
	if err := c.ShouldBindBodyWithJSON(&creatureSpawnForm); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// check if creature exist
	creature, creatureErr := repository.GetCreatureByID(ctx, db, id)
	if creatureErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// check if location exist
	location, locationErr := repository.GetLocationByID(ctx, db, creatureSpawnForm.EmplacementID)
	if locationErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	// check if location do not already assign
	var creatureSpawn entities.CreatureSpawns
	err := pgxscan.Get(ctx, db, &creatureSpawn, `SELECT * FROM creaturespawn WHERE creature_id = $1 AND emplacement_id = $2`, creature.ID, location.ID)
	if errors.Is(pgx.ErrNoRows, err) {
		_, insertErr := db.Exec(ctx, `INSERT INTO creaturespawn (creature_id, emplacement_id, level_required, spawn_rate) values ($1,$2,$3,$4)`,
			creature.ID, location.ID, creatureSpawnForm.LevelRequired, creatureSpawnForm.SpawnRate)
		if insertErr != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}

		c.JSON(http.StatusOK, &creatureSpawn)
		c.Done()
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	if creatureSpawn.CreatureID == 0 {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
}
func UpdateCreatureSkill(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var creatureSkillForm entities.CreatureSkill
	if err := c.ShouldBindBodyWithJSON(&creatureSkillForm); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// check if creature exist
	creature, creatureErr := repository.GetCreatureByID(ctx, db, id)
	if creatureErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// check if skill exist
	skill, skillErr := repository.GetSkillByID(ctx, db, creatureSkillForm.SkillID)
	if skillErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// check if creature do not already have this skill
	var creatureSkill entities.CreatureSkill
	err := pgxscan.Get(ctx, db, &creatureSkill, `SELECT * FROM creature_skill WHERE creature_id = $1 AND skill_id = $2`, creature.ID, skill.ID)
	if errors.Is(pgx.ErrNoRows, err) {
		_, insertErr := db.Exec(ctx, `INSERT INTO creaturespawn (creature_id, emplacement_id, level_required, spawn_rate) values ($1,$2,$3,$4)`,
			creature.ID, skill.ID)
		if insertErr != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}

		c.JSON(http.StatusOK, &creatureSkill)
		c.Done()
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
}
func UpdateCreatureLoot(c *gin.Context) {

}

// DELETE \\

func DeleteCreature(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	// check if creature exist
	creature, err := repository.GetCreatureByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	// delete creature skill
	_, deleteErr := db.Exec(ctx, `DELETE FROM creature_skill WHERE creature_id = $1`, creature.ID)
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	// delete creature spawn
	_, deleteErr = db.Exec(ctx, `DELETE FROM creaturespawn WHERE creature_id = $1`, creature.ID)
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	if creature.IsPet {
		// select pet
		var selectedPets entities.PetsMounts
		err := pgxscan.Get(ctx, db, &selectedPets, `SELECT * FROM pets_mounts where creature_id = $1`, creature.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		// delete pets to all user
		_, deleteErr = db.Exec(ctx, `DELETE FROM user_pets_mounts where pet_id = $1`, selectedPets.ID)
		if deleteErr != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		// delete pet
		_, deleteErr = db.Exec(ctx, `DELETE FROM pets_mounts where creature_id = $1`, creature.ID)
		if deleteErr != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
	}
	// delete creature
	_, deleteErr = db.Exec(ctx, `DELETE FROM creatures where id = $1`, creature.ID)
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.Status(http.StatusOK)
	c.Done()
}
func DeleteCreatureSpawn(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var creatureSpawnForm int
	if err := c.ShouldBindBodyWithJSON(&creatureSpawnForm); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// check if creature exist
	creature, err := repository.GetCreatureByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// delete creature spawn
	_, deleteErr := db.Exec(ctx, `DELETE FROM creaturespawn WHERE creature_id = $1 and emplacement_id = $2`, creature.ID, creatureSpawnForm)
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.Status(http.StatusOK)
	c.Done()
}
func DeleteCreatureSkill(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var creatureSkillForm int
	if err := c.ShouldBindBodyWithJSON(&creatureSkillForm); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// check if creature exist
	creature, err := repository.GetCreatureByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// delete creature spawn
	_, deleteErr := db.Exec(ctx, `DELETE FROM creature_skill WHERE creature_id = $1 and skill_id = $2`, creature.ID, creatureSkillForm)
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.Status(http.StatusOK)
	c.Done()
}
func DeleteCreatureLoot(c *gin.Context) {

}
