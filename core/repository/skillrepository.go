package repository

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"sanctuary-api/entities"
)

func GetSkillByID(ctx context.Context, db *pgxpool.Pool, skillID string) (entities.Skill, error) {
	var skill entities.Skill
	err := pgxscan.Get(ctx, db, &skill, `SELECT id, name, type, description FROM skills where id = $1`, skillID)
	if err != nil {
		return skill, err
	}
	return skill, nil
}
