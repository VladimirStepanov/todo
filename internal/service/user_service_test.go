package service

import (
	"errors"
	"testing"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/VladimirStepanov/todo-app/internal/models/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestConfirmEmail(t *testing.T) {
	retErr := errors.New("Some error")

	tests := []struct {
		name   string
		retErr error
		expErr error
	}{
		{"No error", nil, nil},
		{"With error", retErr, retErr},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoMock := new(mocks.UserRepository)
			repoMock.On("ConfirmEmail", mock.Anything).Return(tc.retErr)
			us := NewUserService(repoMock)

			err := us.ConfirmEmail("test")

			require.Equal(t, tc.expErr, err)
			repoMock.AssertExpectations(t)
		})
	}

}

func TestCreate(t *testing.T) {
	retErr := errors.New("Some error")

	tests := []struct {
		name  string
		err   error
		user  *models.User
		email string
	}{
		{"No error", nil, &models.User{Email: "hello@world.ru"}, "hello@world.ru"},
		{"With error", retErr, nil, "error@world.ru"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoMock := new(mocks.UserRepository)
			mockCall := repoMock.On("Create", mock.Anything)
			mockCall.Run(func(args mock.Arguments) {
				mockCall.Return(args.Get(0), tc.err)
			})
			us := NewUserService(repoMock)

			user, err := us.Create(tc.email, "World")

			require.Equal(t, tc.err, err)
			if err != nil && tc.err == err {
				require.NotEmpty(t, user.ActivatedLink)
				require.Equal(t, tc.email, user.Email)
				require.NotEmpty(t, user.Password)
				require.False(t, user.IsActivated)

			}
			repoMock.AssertExpectations(t)
		})
	}
}
