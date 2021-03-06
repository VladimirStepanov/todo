package service

import (
	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo models.UserRepository
}

func NewUserService(repo models.UserRepository) models.UserService {
	return &UserService{repo: repo}
}

func (us *UserService) Create(Email, Password string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	link := uuid.NewString()

	user := &models.User{
		Email:         Email,
		Password:      string(hashedPassword),
		ActivatedLink: link,
	}

	return us.repo.Create(user)
}

func (us *UserService) ConfirmEmail(Link string) error {
	return us.repo.ConfirmEmail(Link)
}

func (us *UserService) SignIn(Email, Password string) (*models.User, error) {
	user, err := us.repo.FindUserByEmail(Email)

	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Password)) != nil {
		return nil, models.ErrBadUser
	}

	if !user.IsActivated {
		return nil, models.ErrUserNotActivated
	}
	return user, nil
}
