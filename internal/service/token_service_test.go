package service

import (
	"testing"
	"time"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/VladimirStepanov/todo-app/internal/models/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	accessKey   = "accessKey"
	refreshKey  = "refreshKey"
	maxLoggenIn = 6
	testUUID    = "60a1cc8e-f741-45bc-a794-1ac655790c3b"
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

			require.Equal(t, tc.expErr, err)

			if tc.expErr == nil {
				require.NotEmpty(t, td.AccessToken)
				require.NotEmpty(t, td.RefreshToken)
				require.NotEmpty(t, td.UUID)
			}

		})
	}

}

func TestRefresh(t *testing.T) {

	expiredToken, err := GenerateToken(testUUID, userID, 100, 103, refreshKey)
	require.NoError(t, err)
	actualToken, err := GenerateToken(
		testUUID, userID,
		time.Now().Unix(),
		time.Now().Add(time.Hour).Unix(), refreshKey,
	)

	require.NoError(t, err)

	tests := []struct {
		name            string
		token           string
		getRetVal       bool
		getRetErr       error
		delRetErr       error
		countRetVal     int
		countRetErr     error
		expErr          error
		setTokensRetErr error
	}{
		{
			name:            "Bad token",
			token:           "bad.bad.bad",
			getRetVal:       false,
			getRetErr:       nil,
			delRetErr:       nil,
			countRetVal:     0,
			countRetErr:     nil,
			setTokensRetErr: nil,
			expErr:          models.ErrBadToken,
		},
		{
			name:            "Expired token",
			token:           expiredToken,
			getRetVal:       false,
			getRetErr:       nil,
			delRetErr:       nil,
			countRetVal:     0,
			countRetErr:     nil,
			setTokensRetErr: nil,
			expErr:          models.ErrTokenExpired,
		},
		{
			name:            "Get unknown error",
			token:           actualToken,
			getRetVal:       false,
			getRetErr:       ErrSome,
			delRetErr:       nil,
			countRetVal:     0,
			countRetErr:     nil,
			setTokensRetErr: nil,
			expErr:          ErrSome,
		},
		{
			name:            "Get return false",
			token:           actualToken,
			getRetVal:       false,
			getRetErr:       nil,
			delRetErr:       nil,
			countRetVal:     0,
			countRetErr:     nil,
			setTokensRetErr: nil,
			expErr:          models.ErrUserUnauthorized,
		},
		{
			name:            "Delete return unkown error",
			token:           actualToken,
			getRetVal:       true,
			getRetErr:       nil,
			delRetErr:       ErrSome,
			countRetVal:     0,
			countRetErr:     nil,
			setTokensRetErr: nil,
			expErr:          ErrSome,
		},
		{
			name:            "Count return error",
			token:           actualToken,
			getRetVal:       true,
			getRetErr:       nil,
			delRetErr:       nil,
			countRetVal:     0,
			countRetErr:     ErrSome,
			setTokensRetErr: nil,
			expErr:          ErrSome,
		},
		{
			name:            "SetTokens return error",
			token:           actualToken,
			getRetVal:       true,
			getRetErr:       nil,
			delRetErr:       nil,
			countRetVal:     1,
			countRetErr:     nil,
			setTokensRetErr: ErrSome,
			expErr:          ErrSome,
		},
		{
			name:            "Success refresh",
			token:           actualToken,
			getRetVal:       true,
			getRetErr:       nil,
			delRetErr:       nil,
			countRetVal:     1,
			countRetErr:     nil,
			setTokensRetErr: nil,
			expErr:          nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoMock := new(mocks.TokenRepository)
			repoMock.On("Get", mock.Anything).Return(tc.getRetVal, tc.getRetErr)
			repoMock.On("Delete", mock.Anything, mock.Anything).Return(tc.delRetErr)
			repoMock.On("Count", mock.Anything).
				Return(tc.countRetVal, tc.countRetErr)
			repoMock.On("SetTokens", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(tc.setTokensRetErr)

			ts := NewTokenService(accessKey, refreshKey, maxLoggenIn, repoMock)

			td, err := ts.Refresh(tc.token)

			require.Equal(t, tc.expErr, err)

			if tc.expErr == nil {
				require.NotEmpty(t, td.AccessToken)
				require.NotEmpty(t, td.RefreshToken)
				require.NotEmpty(t, td.UUID)
			}

		})
	}
}

func TestVerify(t *testing.T) {
	expiredToken, err := GenerateToken(testUUID, userID, 100, 103, accessKey)
	require.NoError(t, err)
	actualToken, err := GenerateToken(
		testUUID, userID,
		time.Now().Unix(),
		time.Now().Add(time.Hour).Unix(), accessKey,
	)

	require.NoError(t, err)
	tests := []struct {
		name      string
		token     string
		getRetVal bool
		getRetErr error
		expErr    error
	}{
		{
			name:      "Bad token",
			token:     "bad.bad.bad",
			getRetVal: false,
			getRetErr: nil,
			expErr:    models.ErrBadToken,
		},
		{
			name:      "Expired token",
			token:     expiredToken,
			getRetVal: false,
			getRetErr: nil,
			expErr:    models.ErrTokenExpired,
		},
		{
			name:      "Get unknown error",
			token:     actualToken,
			getRetVal: false,
			getRetErr: ErrSome,
			expErr:    ErrSome,
		},
		{
			name:      "Get return false",
			token:     actualToken,
			getRetVal: false,
			getRetErr: nil,
			expErr:    models.ErrUserUnauthorized,
		},
		{
			name:      "Success refresh",
			token:     actualToken,
			getRetVal: true,
			getRetErr: nil,
			expErr:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoMock := new(mocks.TokenRepository)
			repoMock.On("Get", mock.Anything).Return(tc.getRetVal, tc.getRetErr)
			ts := NewTokenService(accessKey, refreshKey, maxLoggenIn, repoMock)

			id, uuid, err := ts.Verify(tc.token)

			require.Equal(t, tc.expErr, err)

			if tc.expErr == nil {
				require.Equal(t, userID, id)
				require.Equal(t, testUUID, uuid)
			}

		})
	}
}

func TestLogout(t *testing.T) {
	tests := []struct {
		name      string
		delRetErr error
		expErr    error
	}{
		{
			name:      "Delete return error",
			delRetErr: ErrSome,
			expErr:    ErrSome,
		},
		{
			name:      "Success logout",
			delRetErr: nil,
			expErr:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoMock := new(mocks.TokenRepository)
			repoMock.On("Delete", mock.Anything, mock.Anything).Return(tc.delRetErr)

			ts := NewTokenService(accessKey, refreshKey, maxLoggenIn, repoMock)
			err := ts.Logout(1, "hello")
			require.Equal(t, tc.expErr, err)
		})
	}
}
