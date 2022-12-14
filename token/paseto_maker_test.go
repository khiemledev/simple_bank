package token

import (
	"testing"
	"time"

	"github.com/khiemledev/simple_bank_golang/util"
	"github.com/stretchr/testify/require"
)

func TestShortPasetoSymmetricKey(t *testing.T) {
	maker, err := CreatePasetoMaker(util.RandomString(minSecretKeySize - 1))
	require.Error(t, err)
	require.Nil(t, maker)
}

func TestPasetoMaker(t *testing.T) {
	maker, err := CreatePasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	duration := time.Minute

	username := util.RandomOwner()
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
}

func TestPasetoTokenExpired(t *testing.T) {
	maker, err := CreatePasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomOwner()
	token, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
