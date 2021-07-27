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
}

type TokenRepository interface {
	Get(key string) (bool, error)
	Set(key string, exp time.Duration) error
	Count(pattern string) (int, error)
}
