package postgres

import (
	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/jmoiron/sqlx"
)

type PostgresListRepository struct {
	DB *sqlx.DB
}

func NewPostgresListRepository(db *sqlx.DB) models.ListRepository {
	return &PostgresListRepository{
		DB: db,
	}
}

func (ls *PostgresListRepository) Create(title, description string, userID int64) (int64, error) {
	return 0, nil
}

func (ls *PostgresListRepository) GrantRole(listID, fromUser, toUserID int64, role bool) error {
	return nil
}

func (ls *PostgresListRepository) GetListByID(listID, userID int64) (*models.List, error) {
	return nil, nil
}

func (ls *PostgresListRepository) GetUserLists(userID int64) ([]*models.List, error) {
	return nil, nil
}

func (ls *PostgresListRepository) Delete(listID, userID int64) error {
	return nil
}

func (ls *PostgresListRepository) Update(userID int64, list *models.List) error {
	return nil
}
