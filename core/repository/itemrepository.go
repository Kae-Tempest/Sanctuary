package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"sanctuary-api/entities"
)

func GetItemByID(ctx context.Context, db *pgxpool.Pool, itemID string) (entities.Item, error) {
	rows := db.QueryRow(ctx, `SELECT id ,name, description, type, rank,
	    strength,constitution,mana,stamina,dexterity,intelligence,wisdom,charisma,enchantment_level,emplacement
		FROM items full join item_stats  on items.id = item_stats.item_id  full join item_emplacement on items.id = item_emplacement.item_id where id = $1`, itemID)

	item, scanErr := assignOneRow(rows)
	if scanErr != nil {
		fmt.Println(scanErr, "scanErr")
		return item, scanErr
	}

	return item, nil
}

func GetItemInfoByName(ctx context.Context, db *pgxpool.Pool, itemName string) (entities.Item, error) {
	var item entities.Item
	err := pgxscan.Get(ctx, db, &item, `SELECT id, name, description, type, rank FROM items where name = $1`, itemName)
	if err != nil {
		return item, err
	}

	return item, nil
}

func GetItemByName(ctx context.Context, db *pgxpool.Pool, itemName string) (entities.Item, error) {
	rows := db.QueryRow(ctx, `SELECT id, name, description, type, rank,
	    strength,constitution,mana,stamina,dexterity,intelligence,wisdom,charisma,enchantment_level,emplacement
		FROM items full join item_stats  on items.id = item_stats.item_id  full join item_emplacement on items.id = item_emplacement.item_id where name = $1`, itemName)

	item, scanErr := assignOneRow(rows)

	if scanErr != nil {
		slog.Error("Error Selection Item (scan return)", slog.Any("error", scanErr))
		return item, scanErr
	}
	return item, nil
}

func GetItemsByType(ctx context.Context, db *pgxpool.Pool, itemType string) ([]entities.Item, error) {

	rows, err := db.Query(ctx, `SELECT name, description, type, rank,
	    strength,constitution,mana,stamina,dexterity,intelligence,wisdom,charisma,enchantment_level,emplacement
		FROM items full join item_stats  on items.id = item_stats.item_id  full join item_emplacement on items.id = item_emplacement.item_id where type = $1`, itemType)
	if err != nil {
		return []entities.Item{}, err
	}

	items, scanErr := AssignMultipleRows(rows)
	if scanErr != nil {
		return items, scanErr
	}

	return items, nil
}

func GetItemStat(ctx context.Context, db *pgxpool.Pool, itemID int) (entities.ItemStat, error) {
	var item entities.ItemStat
	err := pgxscan.Get(ctx, db, &item, `SELECT strength, constitution, mana, stamina, dexterity, intelligence, wisdom, charisma, enchantment_level FROM item_stats where item_id = $1`, itemID)
	if err != nil {
		return item, err
	}

	return item, nil
}

func GetItemEmplacement(ctx context.Context, db *pgxpool.Pool, itemID int) (entities.ItemEmplacement, error) {
	var item entities.ItemEmplacement
	err := pgxscan.Get(ctx, db, &item, `SELECT item_id, emplacement FROM item_emplacement where item_id = $1`, itemID)
	if err != nil {
		return item, err
	}

	return item, nil
}

func assignOneRow(rows pgx.Row) (entities.Item, error) {
	item := entities.Item{}

	var id sql.Null[int]
	var name string
	var description string
	var itemType int
	var rank string
	var strength int
	var constitution int
	var mana int
	var stamina int
	var dexterity int
	var intelligence int
	var wisdom int
	var charisma int
	var enchantmentLevel int
	var emplacement sql.Null[int]

	scanErr := rows.Scan(&id, &name, &description, &itemType, &rank, &strength, &constitution, &mana, &stamina, &dexterity, &intelligence, &wisdom, &charisma, &enchantmentLevel, &emplacement)
	if errors.Is(scanErr, pgx.ErrNoRows) {
		slog.Error("Error Selection Item (scan) : No Rows", slog.Any("error", scanErr))
		return item, scanErr
	}
	if scanErr != nil && !errors.Is(scanErr, pgx.ErrNoRows) {
		slog.Error("Error Selection Item (scan)", slog.Any("error", scanErr))
		return item, scanErr
	}

	item = entities.Item{
		ID:          id.V,
		Name:        name,
		Description: description,
		Type:        itemType,
		Rank:        rank,
		Stats: entities.ItemStat{
			Strength:         strength,
			Constitution:     constitution,
			Mana:             mana,
			Stamina:          stamina,
			Dexterity:        dexterity,
			Intelligence:     intelligence,
			Wisdom:           wisdom,
			Charisma:         charisma,
			EnchantmentLevel: enchantmentLevel,
		},
		ItemEmplacement: entities.ItemEmplacement{
			Emplacement: emplacement.V,
		},
	}
	return item, nil
}

func AssignMultipleRows(rows pgx.Rows) ([]entities.Item, error) {
	var items []entities.Item
	defer rows.Close()

	for rows.Next() {
		var id sql.Null[int]
		var name string
		var description string
		var itemType int
		var rank string
		var strength int
		var constitution int
		var mana int
		var stamina int
		var dexterity int
		var intelligence int
		var wisdom int
		var charisma int
		var enchantmentLevel int
		var emplacement sql.Null[int]

		scanErr := rows.Scan(&id, &name, &description, &itemType, &rank, &strength, &constitution, &mana, &stamina, &dexterity, &intelligence, &wisdom, &charisma, &enchantmentLevel, &emplacement)
		if scanErr != nil {
			return nil, scanErr
		}

		item := entities.Item{
			ID:          id.V,
			Name:        name,
			Description: description,
			Type:        itemType,
			Rank:        rank,
			Stats: entities.ItemStat{
				Strength:         strength,
				Constitution:     constitution,
				Mana:             mana,
				Stamina:          stamina,
				Dexterity:        dexterity,
				Intelligence:     intelligence,
				Wisdom:           wisdom,
				Charisma:         charisma,
				EnchantmentLevel: enchantmentLevel,
			},
			ItemEmplacement: entities.ItemEmplacement{
				Emplacement: emplacement.V,
			},
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
