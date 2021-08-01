package postgres

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var (
	ErrUnknown = errors.New("unknown error")
)

func TestCreate(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Error while sqlmock.New()", err)
	}

	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	lr := NewPostgresListRepository(db)

	tests := []struct {
		name    string
		setMock func(m sqlmock.Sqlmock, e error)
		expErr  error
		expID   int64
	}{
		{
			name: "Begin return error",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectBegin().WillReturnError(e)
			},
			expErr: ErrUnknown,
			expID:  0,
		},
		{
			name: "QueryRow return error",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectBegin()
				m.ExpectQuery("INSERT INTO lists").WithArgs("title", "description").WillReturnError(e)
				m.ExpectRollback()
			},
			expErr: ErrUnknown,
			expID:  0,
		},
		{
			name: "Exec return error",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				m.ExpectQuery("INSERT INTO lists").WithArgs("title", "description").WillReturnRows(rows)
				m.ExpectExec("INSERT INTO users_lists").WithArgs(1, 1, true).WillReturnError(e)
				m.ExpectRollback()
			},
			expErr: ErrUnknown,
			expID:  0,
		},
		{
			name: "Commit return error",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				m.ExpectQuery("INSERT INTO lists").WithArgs("title", "description").WillReturnRows(rows)
				m.ExpectExec("INSERT INTO users_lists").WithArgs(1, 1, true).WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit().WillReturnError(e)
			},
			expErr: ErrUnknown,
			expID:  0,
		},
		{
			name: "Success create",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				m.ExpectQuery("INSERT INTO lists").WithArgs("title", "description").WillReturnRows(rows)
				m.ExpectExec("INSERT INTO users_lists").WithArgs(1, 1, true).WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
			expErr: nil,
			expID:  1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setMock(mock, tc.expErr)
			id, err := lr.Create("title", "description", 1)
			require.Equal(t, tc.expID, id)
			require.Equal(t, tc.expErr, err)
		})
	}
}
