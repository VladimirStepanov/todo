package postgres

import (
	"database/sql"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/jmoiron/sqlx"
)

type PostgresItemRepository struct {
	DB *sqlx.DB
}

func NewPostgresItemRepository(db *sqlx.DB) models.ItemRepository {
	return &PostgresItemRepository{
		DB: db,
	}
}

func (ir *PostgresItemRepository) Create(title, description string, listID int64) (int64, error) {
	var itemID int64

	err := ir.DB.QueryRow(
		`INSERT INTO items(list_id, title, description) VALUES($1, $2, $3) RETURNING id`,
		listID, title, description,
	).Scan(&itemID)

	if err != nil {
		return 0, err
	}

	return itemID, nil
}

func (ir *PostgresItemRepository) GetItems(listID int64) ([]*models.Item, error) {
	return nil, nil
}

func (ir *PostgresItemRepository) GetItemByID(listID, itemID int64) (*models.Item, error) {
	res := &models.Item{}

	err := ir.DB.Get(res, "SELECT * FROM items WHERE list_id=$1 AND id=$2", listID, itemID)

	if err != nil {
		if err == sql.ErrNoRows {
			err = models.ErrNoItem
		}
		return nil, err
	}
	return res, nil
}

func (ir *PostgresItemRepository) Update(listID, itemID int64, item *models.UpdateItemReq) error {
	return nil
}

func (ir *PostgresItemRepository) Delete(listID, itemID int64) error {
	res, err := ir.DB.Exec("DELETE FROM items WHERE id=$1 AND list_id=$2", itemID, listID)

	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if ra == 0 {
		return models.ErrNoItem
	}

	return nil
}
