package models

type Item struct {
	ID          int64  `json:"id" db:"id"`
	ListID      int64  `json:"list_id" db:"list_id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Done        bool   `json:"done"  db:"done"`
}

type UpdateItemReq struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool
}
