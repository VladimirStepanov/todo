package redisrepo

import (
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/require"
)

var ErrUnknown = errors.New("unknown error")

func TestSetTokens(t *testing.T) {
	tests := []struct {
		name    string
		setMock func(m redismock.ClientMock)
		expErr  error
	}{
		{
			name: "Some error",
			setMock: func(m redismock.ClientMock) {
				m.ExpectWatch("test1", "test2").SetErr(ErrUnknown)
			},
			expErr: ErrUnknown,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			tc.setMock(mock)
			rr := NewRedisRepository(db)
			err := rr.SetTokens("test1", time.Hour, "test2", time.Hour)
			require.Equal(t, tc.expErr, err)
		})
	}

}

func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		setMock func(m redismock.ClientMock)
		expVal  bool
		expErr  error
	}{
		{
			name: "Correct get",
			setMock: func(m redismock.ClientMock) {
				m.ExpectGet("hello").SetVal("1")
			},
			expErr: nil,
			expVal: true,
		},
		{
			name: "No value",
			setMock: func(m redismock.ClientMock) {
				m.ExpectGet("hello").RedisNil()
			},
			expErr: nil,
			expVal: false,
		},
		{
			name: "Some error",
			setMock: func(m redismock.ClientMock) {
				m.ExpectGet("hello").SetErr(ErrUnknown)
			},
			expErr: ErrUnknown,
			expVal: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			tc.setMock(mock)
			rr := NewRedisRepository(db)
			val, err := rr.Get("hello")
			require.Equal(t, tc.expErr, err)
			require.Equal(t, tc.expVal, val)
		})
	}

}

func TestCount(t *testing.T) {
	tests := []struct {
		name    string
		setMock func(m redismock.ClientMock)
		expVal  int
		expErr  error
	}{
		{
			name: "Some error",
			setMock: func(m redismock.ClientMock) {
				m.ExpectKeys("pattern").SetErr(ErrUnknown)
			},
			expErr: ErrUnknown,
			expVal: 0,
		},
		{
			name: "No result",
			setMock: func(m redismock.ClientMock) {
				m.ExpectKeys("pattern").SetVal([]string{})
			},
			expErr: nil,
			expVal: 0,
		},
		{
			name: "Some result",
			setMock: func(m redismock.ClientMock) {
				m.ExpectKeys("pattern").SetVal([]string{"1", "2"})
			},
			expErr: nil,
			expVal: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			tc.setMock(mock)
			rr := NewRedisRepository(db)
			val, err := rr.Count("pattern")
			require.Equal(t, tc.expErr, err)
			require.Equal(t, tc.expVal, val)
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name    string
		setMock func(m redismock.ClientMock)
		expErr  error
	}{
		{
			name: "Some error",
			setMock: func(m redismock.ClientMock) {
				m.ExpectDel("hello").SetErr(ErrUnknown)
			},
			expErr: ErrUnknown,
		},
		{
			name: "Success delete",
			setMock: func(m redismock.ClientMock) {
				m.ExpectDel("hello").SetVal(1)
			},
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			tc.setMock(mock)
			rr := NewRedisRepository(db)
			err := rr.Delete("hello")
			require.Equal(t, tc.expErr, err)
		})
	}
}
