package service

import (
	"testing"

	"github.com/VladimirStepanov/todo-app/internal/models/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListCreate(t *testing.T) {
	tests := []struct {
		name   string
		idRet  int64
		errRet error
		expID  int64
		expErr error
	}{
		{
			name:   "Return error",
			idRet:  0,
			errRet: ErrSome,
			expID:  0,
			expErr: ErrSome,
		},
		{
			name:   "Success create",
			idRet:  100,
			errRet: nil,
			expID:  100,
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lr := new(mocks.ListRepository)
			lr.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(tc.idRet, tc.expErr)

			ls := NewListService(lr)

			id, err := ls.Create("title", "description", 1)
			require.Equal(t, tc.expID, id)
			require.Equal(t, tc.expErr, err)
		})
	}
}
