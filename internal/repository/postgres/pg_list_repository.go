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
	tx, err := ls.DB.Beginx()
	if err != nil {
		return 0, err
	}

	var listID int64

	err = tx.QueryRow(
		"INSERT INTO lists(title, description) VALUES($1, $2) RETURNING id",
		title, description,
	).Scan(&listID)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.Exec(
		"INSERT INTO users_lists(user_id, list_id, is_admin) VALUES($1, $2, $3)",
		userID, listID, true,
	)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return listID, nil
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
