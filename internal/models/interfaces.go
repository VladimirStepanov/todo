package models

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
