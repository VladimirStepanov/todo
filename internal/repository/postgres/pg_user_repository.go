package postgres

import (
	"database/sql"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PostgresUserRepository struct {
	DB *sqlx.DB
}

func NewPostgresUserRepository(DB *sqlx.DB) models.UserRepository {
	return &PostgresUserRepository{DB: DB}
}

func (pr *PostgresUserRepository) Create(user *models.User) (*models.User, error) {

	var insertedID int64

	err := pr.DB.QueryRow(
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

func (pr *PostgresUserRepository) ConfirmEmail(Link string) error {
	res, err := pr.DB.Exec("UPDATE users SET is_activated=TRUE WHERE activated_link=$1 AND is_activated=FALSE", Link)
	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return models.ErrConfirmLinkNotExists
	}
	return nil
}

func (pr *PostgresUserRepository) FindUserByEmail(Email string) (*models.User, error) {
	user := &models.User{}

	err := pr.DB.Get(user, "SELECT * FROM users WHERE email=$1", Email)

	if err != nil {
		if err == sql.ErrNoRows {
			err = models.ErrBadUser
		}
		return nil, err
	}

	return user, nil
}
