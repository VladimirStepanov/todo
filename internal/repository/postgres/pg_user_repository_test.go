package postgres

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var testUser = models.User{
	Email:         "test@mail.com",
	Password:      "helloworld",
	IsActivated:   false,
	ActivatedLink: "activated_link",
}

func TestUserCreate(t *testing.T) {
	var retID int64 = 1
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Error while sqlmock.New()", err)
	}

	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	pr := NewPostgresUserRepository(db)

	tests := []struct {
		name    string
		setMock func(m sqlmock.Sqlmock, e error)
		retErr  error
		expErr  error
		expID   int64
	}{
		{
			name: "Return unknown error",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectQuery("INSERT INTO users").
					WithArgs(testUser.Email, testUser.Password, testUser.ActivatedLink).
					WillReturnError(e)
			},
			retErr: ErrUnknown,
			expErr: ErrUnknown,
			expID:  0,
		},
		{
			name: "Return user already exists",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectQuery("INSERT INTO users").
					WithArgs(testUser.Email, testUser.Password, testUser.ActivatedLink).
					WillReturnError(e)
			},
			retErr: &pq.Error{Code: "23505"},
			expErr: models.ErrUserAlreadyExists,
			expID:  0,
		},
		{
			name: "Success create",
			setMock: func(m sqlmock.Sqlmock, e error) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(retID)
				m.ExpectQuery("INSERT INTO users").
					WithArgs(testUser.Email, testUser.Password, testUser.ActivatedLink).
					WillReturnRows(rows)
			},
			retErr: nil,
			expErr: nil,
			expID:  retID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setMock(mock, tc.retErr)
			var inputUser models.User = testUser

			retUser, err := pr.Create(&inputUser)
			require.Equal(t, tc.expErr, err)

			if err == nil {
				require.Equal(t, retID, retUser.ID)
			}
		})
	}
}

func TestConfirmEmail(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Error while sqlmock.New()", err)
	}

	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	pr := NewPostgresUserRepository(db)

	tests := []struct {
		name    string
		setMock func(m sqlmock.Sqlmock, e error)
		retErr  error
		expErr  error
	}{
		{
			name: "Return unknown error",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectExec("UPDATE users").WithArgs("link").WillReturnError(e)
			},
			retErr: ErrUnknown,
			expErr: ErrUnknown,
		},
		{
			name: "Zero rows affected error",
			setMock: func(m sqlmock.Sqlmock, e error) {
				mock.ExpectExec("UPDATE users").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			retErr: nil,
			expErr: models.ErrConfirmLinkNotExists,
		},
		{
			name: "Success confirm",
			setMock: func(m sqlmock.Sqlmock, e error) {
				mock.ExpectExec("UPDATE users").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			retErr: nil,
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setMock(mock, tc.retErr)
			err := pr.ConfirmEmail("link")
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestFindUserByEmail(t *testing.T) {
	var retID int64 = 1
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Error while sqlmock.New()", err)
	}

	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	pr := NewPostgresUserRepository(db)

	tests := []struct {
		name    string
		setMock func(m sqlmock.Sqlmock, e error)
		retErr  error
		expErr  error
		expID   int64
	}{
		{
			name: "Return unknown error",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(testUser.Email).
					WillReturnError(e)
			},
			retErr: ErrUnknown,
			expErr: ErrUnknown,
			expID:  0,
		},
		{
			name: "Return bad user error",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(testUser.Email).
					WillReturnError(e)
			},
			retErr: sql.ErrNoRows,
			expErr: models.ErrBadUser,
			expID:  0,
		},
		{
			name: "Success find",
			setMock: func(m sqlmock.Sqlmock, e error) {
				rows := sqlmock.NewRows(
					[]string{"id", "email", "password_hash", "is_activated", "activated_link"},
				).AddRow(
					retID, testUser.Email,
					testUser.Password,
					testUser.IsActivated,
					testUser.ActivatedLink,
				)

				m.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(testUser.Email).
					WillReturnRows(rows)
			},
			retErr: nil,
			expErr: nil,
			expID:  retID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setMock(mock, tc.retErr)
			var inputUser models.User = testUser
			inputUser.ID = retID

			retUser, err := pr.FindUserByEmail(testUser.Email)
			require.Equal(t, tc.expErr, err)

			if err == nil {
				require.Equal(t, &inputUser, retUser)
			}
		})
	}
}
