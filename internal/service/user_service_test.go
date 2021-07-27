package service

import (
	"errors"
	"testing"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/VladimirStepanov/todo-app/internal/models/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var testEmail = "test@test.ru"
var testPassword = "123456789"
var ErrSome = errors.New("unknown error")

var testUser = &models.User{
	ID:            1,
	Email:         testEmail,
	Password:      "$2a$10$wHVm4AGd.uq.dR7Zk3VjhOJWLEt9WPXEqoCPx5AEzPtH31o7WiY92",
	IsActivated:   false,
	ActivatedLink: "344bda23-9f93-48bf-967c-6b92086baac0",
}

func TestConfirmEmail(t *testing.T) {
	tests := []struct {
		name   string
		retErr error
		expErr error
	}{
		{"No error", nil, nil},
		{"With error", ErrSome, ErrSome},
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
	tests := []struct {
		name  string
		err   error
		user  *models.User
		email string
	}{
		{"No error", nil, &models.User{Email: "hello@world.ru"}, "hello@world.ru"},
		{"With error", ErrSome, nil, "error@world.ru"},
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

func TestSignIn(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		repoRetUser func() *models.User
		repoRetErr  error
		expErr      error
		expUser     func() *models.User
	}{
		{
			name:     "Return unknown error",
			password: testPassword,
			repoRetUser: func() *models.User {
				return nil
			},
			repoRetErr: ErrSome,
			expErr:     ErrSome,
			expUser: func() *models.User {
				return nil
			},
		},
		{
			name:     "User not found",
			password: testPassword,
			repoRetUser: func() *models.User {
				return nil
			},
			repoRetErr: models.ErrBadUser,
			expErr:     models.ErrBadUser,
			expUser: func() *models.User {
				return nil
			},
		},
		{
			name:     "User not found",
			password: "bad_password",
			repoRetUser: func() *models.User {
				return testUser
			},
			repoRetErr: nil,
			expErr:     models.ErrBadUser,
			expUser: func() *models.User {
				return nil
			},
		},
		{
			name:     "User is not activated",
			password: testPassword,
			repoRetUser: func() *models.User {
				return testUser
			},
			repoRetErr: nil,
			expErr:     models.ErrUserNotActivated,
			expUser: func() *models.User {
				return nil
			},
		},
		{
			name:     "Valid user",
			password: testPassword,
			repoRetUser: func() *models.User {
				var u models.User = *testUser
				u.IsActivated = true
				return &u
			},
			repoRetErr: nil,
			expErr:     nil,
			expUser: func() *models.User {
				var u models.User = *testUser
				u.IsActivated = true
				return &u
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoMock := new(mocks.UserRepository)
			repoMock.On("FindUserByEmail", mock.AnythingOfType("string")).Return(tc.repoRetUser(), tc.repoRetErr)
			us := NewUserService(repoMock)

			u, err := us.SignIn(testEmail, tc.password)
			require.Equal(t, err, tc.expErr)
			require.Equal(t, u, tc.expUser())
			repoMock.AssertExpectations(t)
		})
	}
}
