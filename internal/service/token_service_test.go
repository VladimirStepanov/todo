package service

import (
	"testing"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/VladimirStepanov/todo-app/internal/models/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	accessKey   = "accessKey"
	refreshKey  = "refreshKey"
	maxLoggenIn = 6
	userID      = int64(1)
)

func TestNewTokenPair(t *testing.T) {
	tests := []struct {
		name        string
		countRetVal int
		countRetErr error
		setTokRet   error
		expErr      error
	}{
		{
			name:        "Count unknown error",
			countRetVal: 0,
			countRetErr: ErrSome,
			setTokRet:   nil,
			expErr:      ErrSome,
		},
		{
			name:        "Test ErrMaxLoggedIn error",
			countRetVal: 6,
			countRetErr: nil,
			setTokRet:   nil,
			expErr:      models.ErrMaxLoggedIn,
		},
		{
			name:        "SetTokens unknown error",
			countRetVal: 0,
			countRetErr: nil,
			setTokRet:   ErrSome,
			expErr:      ErrSome,
		},
		{
			name:        "Success",
			countRetVal: 0,
			countRetErr: nil,
			setTokRet:   nil,
			expErr:      nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoMock := new(mocks.TokenRepository)
			repoMock.On("Count", mock.Anything).Return(tc.countRetVal, tc.countRetErr)
			repoMock.On("SetTokens", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.setTokRet)

			ts := NewTokenService(accessKey, refreshKey, maxLoggenIn, repoMock)

			td, err := ts.NewTokenPair(userID)

			require.Equal(t, err, tc.expErr)

			if tc.expErr == nil {
				require.NotEmpty(t, td.AccessToken)
				require.NotEmpty(t, td.RefreshToken)
				require.NotEmpty(t, td.UUID)
			}

		})
	}

}
