package repository

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"sanctuary-api/entities"
)

func GetMobByID(ctx context.Context, db *pgxpool.Pool, mobID string) (entities.Mob, error) {
	var mob entities.Mob
	err := pgxscan.Get(ctx, db, &mob, `SELECT id, name, is_pet, level, hp, strength, constitution, mana, stamina, dexterity, intelligence, charisma, wisdom FROM mobs where id = $1`, mobID)
	if err != nil {
		return mob, err
	}
	return mob, nil
}

func GetMobByName(ctx context.Context, db *pgxpool.Pool, mobName string) (entities.Mob, error) {
	var mob entities.Mob
	err := pgxscan.Get(ctx, db, &mob, `SELECT id, name, is_pet, level, hp, strength, constitution, mana, stamina, dexterity, intelligence, charisma, wisdom FROM mobs where name = $1`, mobName)
	if err != nil {
		return mob, err
	}
	return mob, nil
}

func GetMobsByLocation(ctx context.Context, db *pgxpool.Pool, locationID string) ([]entities.Mob, error) {
	var mobs []entities.Mob
	err := pgxscan.Select(ctx, db, &mobs, `SELECT id, name, is_pet, level, hp, strength, constitution, mana, stamina, dexterity, intelligence, charisma, wisdom FROM mobs c left join mob_spawn cp on c.id = cp.mob_id where location_id = $1`, locationID)
	if err != nil {
		return mobs, err
	}
	return mobs, nil
}

func UpdateMobsLocation(ctx context.Context, db *pgxpool.Pool, locationID int, mobs []entities.Mob) error {
	for _, mob := range mobs {
		_, err := db.Exec(ctx, `UPDATE mob_spawn SET location_id = $2 where mob_id = $1`, mob.ID, locationID)
		if err != nil {
			return err
		}
	}
	return nil
}
