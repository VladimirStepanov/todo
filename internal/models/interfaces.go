package models

import "time"

type UserService interface {
	Create(Email, Password string) (*User, error)
	ConfirmEmail(Link string) error
	SignIn(Email, Password string) (*User, error)
}

type UserRepository interface {
	Create(user *User) (*User, error)
	ConfirmEmail(Link string) error
	FindUserByEmail(Email string) (*User, error)
}

type MailService interface {
	SendConfirmationsEmail(user *User) error
}

type TokenService interface {
	NewTokenPair(userID int64) (*TokenDetails, error)
	Refresh(refreshToken string) (*TokenDetails, error)
	Verify(token string) (int64, string, error)
	Logout(userID int64, userUUID string) error
}

type TokenRepository interface {
	Get(key string) (bool, error)
	SetTokens(accessKey string, accessExp time.Duration, refreshKey string, refreshExp time.Duration) error
	Count(pattern string) (int, error)
	Delete(keys ...string) error
}

type ListService interface {
	Create(title, description string, userID int64) (int64, error)
	EditRole(listID, userID int64, role bool) error
	GetListByID(listID, userID int64) (*List, error)
	GetUserLists(userID int64) ([]*List, error)
	Delete(listID int64) error
	Update(list *List) error
	IsListAdmin(ListID, userID int64) error
}

type ListRepository interface {
	Create(title, description string, userID int64) (int64, error)
	EditRole(listID, userID int64, role bool) error
	GetListByID(listID, userID int64) (*List, error)
	GetUserLists(userID int64) ([]*List, error)
	Delete(listID int64) error
	Update(list *List) error
	IsListAdmin(ListID, userID int64) error
}
