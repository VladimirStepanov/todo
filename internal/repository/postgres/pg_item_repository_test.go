package postgres

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var (
	testItem = &models.Item{
		ID:          1,
		ListID:      testList.ID,
		Title:       "title",
		Description: "Description",
		Done:        false,
	}
)

func TestItemCreate(t *testing.T) {
	var retID int64 = 1
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Error while sqlmock.New()", err)
	}

	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	ir := NewPostgresItemRepository(db)

	tests := []struct {
		name    string
		setMock func(m sqlmock.Sqlmock, e error)
		retErr  error
		expErr  error
		expID   int64
	}{
		{
			name: "QueryRow return error",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectQuery("INSERT INTO items").
					WithArgs(1, "title", "description").
					WillReturnError(e)
			},
			retErr: ErrUnknown,
			expErr: ErrUnknown,
			expID:  0,
		},
		{
			name: "Success create",
			setMock: func(m sqlmock.Sqlmock, e error) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(retID)
				m.ExpectQuery("INSERT INTO items").
					WithArgs(1, "title", "description").
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
			id, err := ir.Create("title", "description", 1)
			require.Equal(t, tc.expID, id)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestGetItem(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Error while sqlmock.New()", err)
	}

	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	ir := NewPostgresItemRepository(db)

	tests := []struct {
		name    string
		setMock func(m sqlmock.Sqlmock, e error)
		retErr  error
		expErr  error
		expItem *models.Item
	}{
		{
			name: "Return unknown error",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectQuery(regexp.QuoteMeta("SELECT * FROM items")).
					WithArgs(testList.ID, 1).
					WillReturnError(e)
			},
			retErr:  ErrUnknown,
			expErr:  ErrUnknown,
			expItem: nil,
		},
		{
			name: "Item not found",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectQuery(regexp.QuoteMeta("SELECT * FROM items")).
					WithArgs(testList.ID, 1).
					WillReturnError(e)
			},
			retErr:  sql.ErrNoRows,
			expErr:  models.ErrNoItem,
			expItem: nil,
		},
		{
			name: "Success get",
			setMock: func(m sqlmock.Sqlmock, e error) {
				rows := sqlmock.NewRows(
					[]string{"id", "list_id", "title", "description", "done"},
				).AddRow(testItem.ID, testItem.ID, testItem.Title, testItem.Description, testItem.Done)
				m.ExpectQuery(regexp.QuoteMeta("SELECT * FROM items")).
					WithArgs(testList.ID, testItem.ID).
					WillReturnRows(rows)
			},
			retErr:  nil,
			expErr:  nil,
			expItem: testItem,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setMock(mock, tc.retErr)

			retItem, err := ir.GetItemByID(testList.ID, 1)
			require.Equal(t, tc.expErr, err)

			if err == nil {
				require.Equal(t, tc.expItem, retItem)
			}
		})
	}
}

func TestDeleteItem(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Error while sqlmock.New()", err)
	}

	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	ir := NewPostgresItemRepository(db)

	tests := []struct {
		name    string
		setMock func(m sqlmock.Sqlmock, e error)
		retErr  error
		expErr  error
	}{
		{
			name: "Delete unknown error",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectExec("DELETE FROM items").
					WithArgs(1, 1).
					WillReturnError(e)
			},
			retErr: ErrUnknown,
			expErr: ErrUnknown,
		},
		{
			name: "Delete return ErrNoItem",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectExec("DELETE FROM items").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			retErr: nil,
			expErr: models.ErrNoItem,
		},
		{
			name: "Success delete",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectExec("DELETE FROM items").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			retErr: nil,
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setMock(mock, tc.retErr)

			err := ir.Delete(1, 1)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestUpdateItem(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal("Error while sqlmock.New()", err)
	}

	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	ir := NewPostgresItemRepository(db)

	req := &models.UpdateItemReq{
		Title:       new(string),
		Description: new(string),
		Done:        new(bool),
	}
	*req.Title = "hello"
	*req.Description = "world"
	*req.Done = true

	tests := []struct {
		name    string
		setMock func(m sqlmock.Sqlmock, e error)
		retErr  error
		expErr  error
	}{
		{
			name: "Update unknown error",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectExec("UPDATE items").
					WithArgs(*req.Title, *req.Description, *req.Done, 1, 1).
					WillReturnError(e)
			},
			retErr: ErrUnknown,
			expErr: ErrUnknown,
		},
		{
			name: "Update return ErrNoList",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectExec("UPDATE items").
					WithArgs(*req.Title, *req.Description, *req.Done, 1, 1).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			retErr: nil,
			expErr: models.ErrNoList,
		},
		{
			name: "Success update",
			setMock: func(m sqlmock.Sqlmock, e error) {
				m.ExpectExec("UPDATE items").
					WithArgs(*req.Title, *req.Description, *req.Done, 1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			retErr: nil,
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setMock(mock, tc.retErr)
			err := ir.Update(1, 1, req)
			require.Equal(t, tc.expErr, err)
		})
	}
}
