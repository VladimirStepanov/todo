package postgres

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var (
	ErrUnknown = errors.New("unknown error")
	testList   = &models.List{
		ID:          1,
		Title:       "hello",
		Description: "world",
	}
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
		retErr  error
		expErr  error
		expID   int64
	}{
		{
			name: "Begin return error",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectBegin().WillReturnError(e)
			},
			retErr: ErrUnknown,
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
			retErr: ErrUnknown,
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
			retErr: ErrUnknown,
			expErr: ErrUnknown,
			expID:  0,
		},
		{
			name: "Commit return error",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				m.ExpectQuery("INSERT INTO lists").
					WithArgs("title", "description").
					WillReturnRows(rows)
				m.ExpectExec("INSERT INTO users_lists").
					WithArgs(1, 1, true).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit().WillReturnError(e)
			},
			retErr: ErrUnknown,
			expErr: ErrUnknown,
			expID:  0,
		},
		{
			name: "Success create",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				m.ExpectQuery("INSERT INTO lists").
					WithArgs("title", "description").
					WillReturnRows(rows)
				m.ExpectExec("INSERT INTO users_lists").
					WithArgs(1, 1, true).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
			retErr: nil,
			expErr: nil,
			expID:  1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setMock(mock, tc.retErr)
			id, err := lr.Create("title", "description", 1)
			require.Equal(t, tc.expID, id)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestGetListByID(t *testing.T) {
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
		retErr  error
		expErr  error
		expList *models.List
	}{
		{
			name: "Return unknown error",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectQuery("SELECT id, title, description FROM lists").
					WithArgs(testList.ID, 1).
					WillReturnError(e)
			},
			retErr:  ErrUnknown,
			expErr:  ErrUnknown,
			expList: nil,
		},
		{
			name: "List not found",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectQuery("SELECT id, title, description FROM lists").
					WithArgs(testList.ID, 1).
					WillReturnError(e)
			},
			retErr:  sql.ErrNoRows,
			expErr:  models.ErrNoList,
			expList: nil,
		},
		{
			name: "Success get",
			setMock: func(m sqlmock.Sqlmock, e error) {
				rows := sqlmock.NewRows(
					[]string{"id", "title", "description"},
				).AddRow(testList.ID, testList.Title, testList.Description)
				m.ExpectQuery("SELECT id, title, description FROM lists").
					WithArgs(testList.ID, 1).
					WillReturnRows(rows)
			},
			retErr:  nil,
			expErr:  nil,
			expList: testList,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setMock(mock, tc.retErr)

			retList, err := lr.GetListByID(testList.ID, 1)
			require.Equal(t, tc.expErr, err)

			if err == nil {
				require.Equal(t, testList, retList)
			}
		})
	}
}