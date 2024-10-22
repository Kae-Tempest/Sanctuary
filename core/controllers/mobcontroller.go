package controllers

import (
	"errors"
	"github.com/jackc/pgx/v5"
	"net/http"
	"sanctuary-api/database"
	"sanctuary-api/entities"
	"sanctuary-api/repository"
	"strconv"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
)

// GET \\

func GetAllMobs(c *gin.Context) {
	db := database.Connect()
	var mobs []entities.Mob
	err := pgxscan.Select(ctx, db, &mobs, `SELECT id, name, is_pet, strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma, level, hp FROM mobs`)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	if len(mobs) > 0 {
		c.JSON(http.StatusOK, &mobs)
	} else {
		c.JSON(http.StatusNotFound, gin.H{})
	}
}
func GetMobByID(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	mob, err := repository.GetMobByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &mob)
}
func GetMobSpawn(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	mob, err := repository.GetMobByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var creatureSpawn entities.MobSpawns
	err = pgxscan.Get(ctx, db, &creatureSpawn, `SELECT mob_id, location_id, level_required, spawn_rate FROM mob_spawn where mob_id = $1`, mob.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &creatureSpawn)
}
func GetMobSkill(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	mob, err := repository.GetMobByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	var creatureSkill []entities.MobSkill
	err = pgxscan.Select(ctx, db, &creatureSkill, `SELECT mob_id, skill_id FROM mob_skill where mob_id = $1`, mob.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &creatureSkill)
}

// POST \\

func CreateMob(c *gin.Context) {
	db := database.Connect()
	var mobForm entities.Mob
	if err := c.ShouldBindBodyWithJSON(&mobForm); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	// check if creature exist
	existingMob, ecErr := repository.GetMobByName(ctx, db, mobForm.Name)
	if errors.Is(ecErr, pgx.ErrNoRows) {
		c.JSON(http.StatusConflict, &existingMob)
	}
	if ecErr != nil && !errors.Is(ecErr, pgx.ErrNoRows) {
		c.String(http.StatusBadRequest, "bad request")
	}
	// insert creature
	_, err := db.Exec(ctx, `INSERT INTO mobs (name, is_pet, strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma, level, hp) values ($1, $2, $3, $4,$5, $6, $7, $8,$9, $10, $11, $12)`,
		mobForm.Name, mobForm.IsPet, mobForm.Strength, mobForm.Constitution, mobForm.Mana, mobForm.Stamina,
		mobForm.Dexterity, mobForm.Intelligence, mobForm.Wisdom, mobForm.Charisma, mobForm.Level, mobForm.HP)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	// Get Creature
	mob, cErr := repository.GetMobByName(ctx, db, mobForm.Name)
	if cErr != nil {
		c.String(http.StatusBadRequest, "bad request")
	}
	// return Creature
	c.JSON(http.StatusCreated, &mob)

}
func AddMobSpawn(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	var mobSpawnForm entities.MobSpawns
	if err := c.ShouldBindBodyWithJSON(&mobSpawnForm); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	mob, err := repository.GetMobByID(ctx, db, id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	var selectedEmplacement entities.Locations
	err = pgxscan.Get(ctx, db, &selectedEmplacement, `SELECT id FROM locations where id = $1`, mobSpawnForm.LocationID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	_, err = db.Exec(ctx, `INSERT INTO mob_spawn (mob_id, location_id, level_required, spawn_rate) values ($1, $2, $3, $4)`,
		mob.ID, selectedEmplacement.ID, mobSpawnForm.LevelRequired, mobSpawnForm.SpawnRate)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	var mobSpawn entities.MobSpawns
	err = pgxscan.Get(ctx, db, &mobSpawn, `SELECT mob_id, location_id, level_required, spawn_rate FROM mob_spawn where mob_id = $1 AND location_id = $2`, mob.ID, selectedEmplacement.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusCreated, &mobSpawn)

}
func AddMobSkill(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")
	var mobSkillForm entities.MobSkill
	if err := c.ShouldBindBodyWithJSON(&mobSkillForm); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	mob, err := repository.GetMobByID(ctx, db, id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	var seletecdSkill entities.Locations
	err = pgxscan.Get(ctx, db, &seletecdSkill, `SELECT id FROM skills where id = $1`, mobSkillForm.SkillID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	_, err = db.Exec(ctx, `INSERT INTO mob_skill (mob_id, skill_id) values ($1,$2)`, id, mobSkillForm.SkillID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	var mobSkill entities.MobSpawns
	err = pgxscan.Get(ctx, db, &mobSkill, `SELECT id, name, description, type, mob_id, skill_id FROM skills s join mob_skill cs on s.id = cs.skill_id where mob_id = $1 and s.id = $2`, mob.ID, seletecdSkill.ID)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusCreated, &mobSkill)
}
func AddMobLoot(c *gin.Context) {}

// PATCH \\

func UpdateMob(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var mobUpdateForm entities.Mob
	if err := c.ShouldBindBodyWithJSON(&mobUpdateForm); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	mob, err := repository.GetMobByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	_, err = db.Exec(ctx, `UPDATE creatures set (name, is_pet, strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma, level, hp ) 
    = ($2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) where id = $1`,
		mob.ID, mobUpdateForm.Name, mobUpdateForm.IsPet, mobUpdateForm.Strength, mobUpdateForm.Mana, mobUpdateForm.Stamina,
		mobUpdateForm.Dexterity, mobUpdateForm.Intelligence, mobUpdateForm.Wisdom, mobUpdateForm.Level, mobUpdateForm.HP)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	mob, err = repository.GetMobByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &mob)

}
func UpdateMobSpawn(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var mobSpawnsForm entities.MobSpawns
	if err := c.ShouldBindBodyWithJSON(&mobSpawnsForm); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// check if creature exist
	creature, creatureErr := repository.GetMobByID(ctx, db, id)
	if creatureErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// check if location exist
	location, locationErr := repository.GetLocationByID(ctx, db, strconv.Itoa(mobSpawnsForm.LocationID))
	if locationErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	// check if location do not already assign
	var mobSpawns entities.MobSpawns
	err := pgxscan.Get(ctx, db, &mobSpawns, `SELECT mob_id, location_id, level_required, spawn_rate FROM mob_spawn WHERE mob_id = $1 AND location_id = $2`, creature.ID, location.ID)
	if errors.Is(pgx.ErrNoRows, err) {
		_, insertErr := db.Exec(ctx, `INSERT INTO mob_spawn (mob_id, location_id, level_required, spawn_rate) values ($1,$2,$3,$4)`,
			creature.ID, location.ID, mobSpawnsForm.LevelRequired, mobSpawnsForm.SpawnRate)
		if insertErr != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}

		c.JSON(http.StatusOK, &mobSpawns)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	if mobSpawns.MobID == 0 {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
}
func UpdateMobSkill(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var mobSkillForm entities.MobSkill
	if err := c.ShouldBindBodyWithJSON(&mobSkillForm); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// check if creature exist
	creature, creatureErr := repository.GetMobByID(ctx, db, id)
	if creatureErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// check if skill exist
	skill, skillErr := repository.GetSkillInfoByID(ctx, db, strconv.Itoa(mobSkillForm.SkillID))
	if skillErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// check if creature do not already have this skill
	var mobSkill entities.MobSkill
	err := pgxscan.Get(ctx, db, &mobSkill, `SELECT mob_id, skill_id FROM mob_skill WHERE mob_id = $1 AND skill_id = $2`, creature.ID, skill.ID)
	if errors.Is(pgx.ErrNoRows, err) {
		_, insertErr := db.Exec(ctx, `INSERT INTO mob_skill (mob_id, skill_id) values ($1,$2)`,
			creature.ID, skill.ID)
		if insertErr != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}

		c.JSON(http.StatusOK, &mobSkill)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
}
func UpdateMobLoot(c *gin.Context) {}

// DELETE \\

func DeleteMob(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	// check if mob exist
	mob, err := repository.GetMobByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	// delete mob skill
	_, deleteErr := db.Exec(ctx, `DELETE FROM mob_skill WHERE mob_id = $1`, mob.ID)
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}
	// delete mob spawn
	_, deleteErr = db.Exec(ctx, `DELETE FROM mob_spawn WHERE mob_id = $1`, mob.ID)
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	if mob.IsPet {
		// select pet
		var selectedPets entities.PetsMounts
		err := pgxscan.Get(ctx, db, &selectedPets, `SELECT mob_id, is_mountable, speed, id FROM pets_mounts where mob_id = $1`, mob.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		// delete pets to all user
		_, deleteErr = db.Exec(ctx, `DELETE FROM character_pets_mounts where pet_id = $1`, selectedPets.ID)
		if deleteErr != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		// delete pet
		_, deleteErr = db.Exec(ctx, `DELETE FROM pets_mounts where mob_id = $1`, mob.ID)
		if deleteErr != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
	}
	// delete mob
	_, deleteErr = db.Exec(ctx, `DELETE FROM mobs where id = $1`, mob.ID)
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.Status(http.StatusOK)
}
func DeleteMobSpawn(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var mobSpawnForm int
	if err := c.ShouldBindBodyWithJSON(&mobSpawnForm); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// check if mob exist
	mob, err := repository.GetMobByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// delete mob spawn
	_, deleteErr := db.Exec(ctx, `DELETE FROM mob_spawn WHERE mob_id = $1 and location_id = $2`, mob.ID, mobSpawnForm)
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.Status(http.StatusOK)
}
func DeleteMobSkill(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var creatureSkillForm int
	if err := c.ShouldBindBodyWithJSON(&creatureSkillForm); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// check if creature exist
	creature, err := repository.GetMobByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	// delete creature spawn
	_, deleteErr := db.Exec(ctx, `DELETE FROM mob_skill WHERE mob_id = $1 and skill_id = $2`, creature.ID, creatureSkillForm)
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.Status(http.StatusOK)
}
func DeleteMobLoot(c *gin.Context) {}
