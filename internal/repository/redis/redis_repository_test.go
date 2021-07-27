package redis

import (
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/require"
)

var ErrUnknown = errors.New("unknown error")

func TestSet(t *testing.T) {
	tests := []struct {
		name    string
		setMock func(m redismock.ClientMock)
		expErr  error
	}{
		{
			name: "Some error",
			setMock: func(m redismock.ClientMock) {
				m.ExpectSet("hello", true, time.Second*10).SetErr(ErrUnknown)
			},
			expErr: ErrUnknown,
		},
		{
			name: "Correct set",
			setMock: func(m redismock.ClientMock) {
				m.ExpectSet("hello", true, time.Second*10).SetVal("OK")
			},
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			tc.setMock(mock)
			rr := NewRedisRepository(db)
			err := rr.Set("hello", time.Second*10)
			require.Equal(t, err, tc.expErr)
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
			require.Equal(t, err, tc.expErr)
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
			require.Equal(t, err, tc.expErr)
			require.Equal(t, tc.expVal, val)
		})
	}
}
