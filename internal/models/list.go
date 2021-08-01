package models

type List struct {
	ID          int64  `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
}

type UsersList struct {
	UserID  int64 `db:"user_id"`
	ListID  int64 `db:"list_id"`
	IsAdmin bool  `db:"is_admin"`
}
