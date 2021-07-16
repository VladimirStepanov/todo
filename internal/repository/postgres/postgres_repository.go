package postgres

import (
	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PostgresRepository struct {
	DB *sqlx.DB
}

func NewPostgresRepository(DB *sqlx.DB) models.UserRepository {
	return &PostgresRepository{DB: DB}
}

func (pr *PostgresRepository) Create(user *models.User) (*models.User, error) {

	var insertedID int64

	err := pr.DB.QueryRowx(
		"INSERT INTO users(email, password_hash, activated_link) values($1, $2, $3) RETURNING id",
		user.Email, user.Password, user.ActivatedLink,
	).Scan(&insertedID)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return nil, models.ErrUserAlreadyExists
			}
		}
		return nil, err
	}
	user.ID = insertedID
	return user, nil
}
