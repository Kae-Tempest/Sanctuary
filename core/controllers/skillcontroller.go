package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"net/http"
	"sanctuary-api/database"
	"sanctuary-api/entities"
	"sanctuary-api/repository"
	"strconv"
)

func GetSkills(c *gin.Context) {
	db := database.Connect()
	rows, err := db.Query(ctx, `SELECT id, name, description, type FROM skills`)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	skills, err := repository.AssignMultipleRowsSkill(rows)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, skills)
}
func GetSkillByID(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	skill, err := repository.GetSkillByID(ctx, db, id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, skill)
}
func GetSkillByType(c *gin.Context) {
	db := database.Connect()
	skillType := c.Param("type")
	rows := db.QueryRow(ctx, `SELECT id, name, description, type FROM skills where type = $1`, skillType)

	skill, err := repository.AssignOneRowSkill(rows)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, skill)
}

func CreateSkill(c *gin.Context) {
	db := database.Connect()
	var skillForm entities.Skill
	if err := c.ShouldBindBodyWithJSON(&skillForm); err != nil {
		slog.Error("Error during binding SkillForm", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	existingSkill, err := repository.GetSkillByName(ctx, db, skillForm.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		_, insertErr := db.Exec(ctx, `INSERT INTO skills (name, description, type) VALUES ($1,$2,$3)`, skillForm.Type, skillForm.Description, skillForm.Type)
		if insertErr != nil {

		}
		skill, getErr := repository.GetSkillByName(ctx, db, skillForm.Name)
		if getErr != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		_, insertErr = db.Exec(ctx, `INSERT INTO skill_stats (skill_id, strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
			skill.ID, skillForm.SkillStat.Strength, skillForm.SkillStat.Constitution, skillForm.SkillStat.Mana, skillForm.SkillStat.Stamina,
			skillForm.SkillStat.Dexterity, skillForm.SkillStat.Intelligence, skillForm.SkillStat.Wisdom, skillForm.SkillStat.Charisma)
		if insertErr != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		newSkill, selectErr := repository.GetSkillByID(ctx, db, strconv.Itoa(skill.ID))
		if selectErr != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}
		c.JSON(http.StatusCreated, &newSkill)
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("Error Getting Item by Name", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error Getting Item by Name")
		return
	}
	if existingSkill.ID != 0 {
		c.String(http.StatusConflict, "already exist")
		return
	}

}

func UpdateSkillInfo(c *gin.Context) {
	db := database.Connect()
	var skillForm entities.Skill
	if err := c.ShouldBindBodyWithJSON(&skillForm); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	skill, err := repository.GetSkillInfoByID(ctx, db, strconv.Itoa(skillForm.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	_, updateErr := db.Exec(ctx, `UPDATE skills set (name, description, type) = ($2,$3,$4) where id = $1`, skill.ID, skillForm.Name, skillForm.Description, skillForm.Type)
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	updatedSkill, gErr := repository.GetSkillInfoByID(ctx, db, strconv.Itoa(skillForm.ID))
	if gErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &updatedSkill)
}

func UpdateSkillStats(c *gin.Context) {
	db := database.Connect()
	var skillForm entities.Skill
	if err := c.ShouldBindBodyWithJSON(&skillForm); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	skill, err := repository.GetSkillByID(ctx, db, strconv.Itoa(skillForm.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	_, updateErr := db.Exec(ctx, `UPDATE skill_stats set (strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma) = ($2,$3,$4,$5,$6,$7,$8,$9) where skill_id = $1`,
		skill.ID, skillForm.SkillStat.Strength, skillForm.SkillStat.Constitution, skillForm.SkillStat.Mana, skillForm.SkillStat.Stamina,
		skillForm.SkillStat.Dexterity, skillForm.SkillStat.Intelligence, skillForm.SkillStat.Wisdom, skillForm.SkillStat.Charisma)
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	updatedSkill, gErr := repository.GetSkillByID(ctx, db, strconv.Itoa(skillForm.ID))
	if gErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusOK, &updatedSkill)
}

func DeleteSkill(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	skill, err := repository.GetSkillByID(ctx, db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	_, deleteErr := db.Exec(ctx, `DELETE FROM skills where id = $1`, skill.ID)
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	_, deleteErr = db.Exec(ctx, `DELETE FROM skill_stats where skill_id = $1`, skill.ID)
	if deleteErr != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	c.Status(http.StatusOK)
}
