package models

type UserService interface {
	Create(Email, Password string) (*User, error)
	ConfirmEmail(Link string) error
}

type UserRepository interface {
	Create(user *User) (*User, error)
	ConfirmEmail(Link string) error
}

type MailService interface {
	SendConfirmationsEmail(user *User) error
}
