package postgres

import (
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

func (is *PostgresItemRepository) Create(title, description string, listID int64) (int64, error) {
	return 0, nil
}

func (is *PostgresItemRepository) GetItems(listID int64) ([]*models.Item, error) {
	return nil, nil
}

func (is *PostgresItemRepository) GetItemBydID(listID, itemID int64) (*models.Item, error) {
	return nil, nil
}

func (is *PostgresItemRepository) Update(listID, itemID int64, item *models.UpdateItemReq) error {
	return nil
}

func (is *PostgresItemRepository) Delete(listID, itemID int64) error {
	return nil
}
