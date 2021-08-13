package models

type List struct {
	ID          int64  `json:"list_id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
}

type UpdateListReq struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type UsersList struct {
	UserID  int64 `db:"user_id"`
	ListID  int64 `db:"list_id"`
	IsAdmin bool  `db:"is_admin"`
}
