package postgres

import (
	"database/sql"
	"fmt"
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

func TestCreateSuccess(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Error while sqlmock.New()", err)
	}

	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	pr := NewPostgresRepository(db)
	var retID int64 = 1
	rows := sqlmock.NewRows([]string{"id"}).AddRow(retID)
	mock.ExpectQuery("INSERT INTO users").
		WithArgs(testUser.Email, testUser.Password, testUser.ActivatedLink).
		WillReturnRows(rows)

	var inputUser models.User = testUser

	retUser, err := pr.Create(&inputUser)

	require.Equal(t, retID, retUser.ID)
	require.NoError(t, err)

}

func TestCreateErrors(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Error while sqlmock.New()", err)
	}

	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	pr := NewPostgresRepository(db)

	unknownError := fmt.Errorf("Unknown error")

	tests := []struct {
		willRetErr error
		expRetErr  error
	}{
		{unknownError, unknownError},
		{&pq.Error{Code: "23505"}, models.ErrUserAlreadyExists},
	}

	for _, tc := range tests {
		mock.ExpectQuery("INSERT INTO users").
			WithArgs(testUser.Email, testUser.Password, testUser.ActivatedLink).
			WillReturnError(tc.willRetErr)
		var inputUser models.User = testUser
		_, err := pr.Create(&inputUser)
		require.EqualError(t, err, tc.expRetErr.Error())
	}
}

func TestConfirmEmail(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Error while sqlmock.New()", err)
	}

	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	tests := []struct {
		name         string
		rowsAffected int
		retErr       error
	}{
		{"Success update", 1, nil},
		{"Error update", 0, models.ErrConfirmLinkNotExists},
	}

	for _, tc := range tests {
		pr := NewPostgresRepository(db)
		mock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(0, int64(tc.rowsAffected)))
		err := pr.ConfirmEmail("testlink")

		require.Equal(t, err, tc.retErr)
	}
}

func TestFindUserByEmailErrors(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Error while sqlmock.New()", err)
	}

	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	pr := NewPostgresRepository(db)

	unknownError := fmt.Errorf("Unknown error")
	tests := []struct {
		willRetErr error
		expRetErr  error
	}{
		{unknownError, unknownError},
		{sql.ErrNoRows, models.ErrBadUser},
	}

	for _, tc := range tests {
		mock.ExpectQuery("SELECT (.+) FROM users").
			WithArgs(testUser.Email).
			WillReturnError(tc.willRetErr)
		_, err := pr.FindUserByEmail(testUser.Email)
		require.EqualError(t, err, tc.expRetErr.Error())
	}
}

func TestFindUserByEmailSuccess(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Error while sqlmock.New()", err)
	}

	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	pr := NewPostgresRepository(db)
	var retID int64 = 1
	rows := sqlmock.NewRows(
		[]string{"id", "email", "password_hash", "is_activated", "activated_link"},
	).AddRow(retID, testUser.Email, testUser.Password, testUser.IsActivated, testUser.ActivatedLink)
	mock.ExpectQuery("SELECT (.+) FROM users").
		WithArgs(testUser.Email).
		WillReturnRows(rows)

	var inputUser models.User = testUser
	inputUser.ID = retID

	retUser, err := pr.FindUserByEmail(testUser.Email)

	require.Equal(t, &inputUser, retUser)
	require.NoError(t, err)

}