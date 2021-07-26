package models

type User struct {
	ID            int64  `db:"id"`
	Email         string `db:"email"`
	Password      string `db:"password_hash"`
	IsActivated   bool   `db:"is_activated"`
	ActivatedLink string `db:"activated_link"`
}
