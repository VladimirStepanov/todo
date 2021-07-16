package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB(pgHost, pgPort, pgUser, pgPassword, pgDB, pgSSL string) (*sqlx.DB, error) {
	pgConnString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pgHost, pgPort, pgUser, pgPassword, pgDB, pgSSL,
	)

	return sqlx.Connect("postgres", pgConnString)
}
