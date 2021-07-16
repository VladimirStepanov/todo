package models

type UserService interface {
	Create(Email, Password string) (*User, error)
}

type UserRepository interface {
	Create(user *User) (*User, error)
}

type MailService interface {
	SendConfirmationsEmail(user *User) error
}
