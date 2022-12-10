package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(t *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(t.Run())
}
