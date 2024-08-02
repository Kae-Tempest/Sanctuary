package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"sanctuary-api/entities"
)

func GetSkillInfoByID(ctx context.Context, db *pgxpool.Pool, skillID string) (entities.Skill, error) {
	var skill entities.Skill
	err := pgxscan.Get(ctx, db, &skill, `SELECT id, name, type, description FROM skills where id = $1`, skillID)
	if err != nil {
		return skill, err
	}
	return skill, nil
}

func GetSkillByName(ctx context.Context, db *pgxpool.Pool, skillName string) (entities.Skill, error) {
	var skill entities.Skill
	err := pgxscan.Get(ctx, db, &skill, `SELECT id, name, type, description FROM skills where id = $1`, skillName)
	if err != nil {
		return skill, err
	}
	return skill, nil
}

func GetSkillByID(ctx context.Context, db *pgxpool.Pool, skillID string) (entities.Skill, error) {
	rows := db.QueryRow(ctx, `SELECT id, name, description, type, strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma 
	FROM skills s join skill_stats on s.id = skill_stats.skill_id where id = $1`, skillID)

	skill, err := AssignOneRowSkill(rows)
	if err != nil {
		return entities.Skill{}, err
	}
	return skill, nil
}

func AssignOneRowSkill(rows pgx.Row) (entities.Skill, error) {
	skill := entities.Skill{}

	var id sql.Null[int]
	var name string
	var description string
	var skillType string
	var strength int
	var constitution int
	var mana int
	var stamina int
	var dexterity int
	var intelligence int
	var wisdom int
	var charisma int

	scanErr := rows.Scan(&id, &name, &description, &skillType, &strength, &constitution, &mana, &stamina, &dexterity, &intelligence, &wisdom, &charisma)
	if errors.Is(scanErr, pgx.ErrNoRows) {
		slog.Error("Error Selection Skill (scan) : No Rows", slog.Any("error", scanErr))
		return skill, scanErr
	}
	if scanErr != nil && !errors.Is(scanErr, pgx.ErrNoRows) {
		slog.Error("Error Selection Skill (scan)", slog.Any("error", scanErr))
		return skill, scanErr
	}

	skill = entities.Skill{
		ID:          id.V,
		Name:        name,
		Description: description,
		Type:        skillType,
		SkillStat: entities.SkillStat{
			Strength:     strength,
			Constitution: constitution,
			Mana:         mana,
			Stamina:      stamina,
			Dexterity:    dexterity,
			Intelligence: intelligence,
			Wisdom:       wisdom,
			Charisma:     charisma,
		},
	}
	return skill, nil
}

func AssignMultipleRowsSkill(rows pgx.Rows) ([]entities.Skill, error) {
	var skills []entities.Skill
	defer rows.Close()

	for rows.Next() {
		var id sql.Null[int]
		var name string
		var description string
		var skillType string
		var strength int
		var constitution int
		var mana int
		var stamina int
		var dexterity int
		var intelligence int
		var wisdom int
		var charisma int

		scanErr := rows.Scan(&id, &name, &description, &skillType, &strength, &constitution, &mana, &stamina, &dexterity, &intelligence, &wisdom, &charisma)
		if scanErr != nil {
			return nil, scanErr
		}

		skill := entities.Skill{
			ID:          id.V,
			Name:        name,
			Description: description,
			Type:        skillType,
			SkillStat: entities.SkillStat{
				Strength:     strength,
				Constitution: constitution,
				Mana:         mana,
				Stamina:      stamina,
				Dexterity:    dexterity,
				Intelligence: intelligence,
				Wisdom:       wisdom,
				Charisma:     charisma,
			},
		}

		skills = append(skills, skill)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return skills, nil
}
