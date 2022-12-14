package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/khiemledev/simple_bank_golang/util"
	"github.com/stretchr/testify/require"
)

func TestShortSecretKey(t *testing.T) {
	maker, err := CreateJWTMaker(util.RandomString(minSecretKeySize - 1))
	require.Error(t, err)
	require.Nil(t, maker)
}

func TestJWTMaker(t *testing.T) {
	maker, err := CreateJWTMaker(util.RandomString(32))
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

func TestTokenExpired(t *testing.T) {
	maker, err := CreateJWTMaker(util.RandomString(32))
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

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	maker, err := CreateJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
