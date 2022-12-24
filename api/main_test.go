package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/khiemledev/simple_bank_golang/db/sqlc"
	"github.com/khiemledev/simple_bank_golang/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store, config *util.Config) *Server {
	if config == nil {
		config = &util.Config{
			TokenSymmetricKey:   util.RandomString(32),
			AccessTokenDuration: time.Minute,
		}
	}
	server, err := NewServer(*config, store)
	require.NoError(t, err)

	return server
}

func TestMain(t *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(t.Run())
}
