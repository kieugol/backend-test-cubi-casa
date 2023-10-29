package tests

import (
	"flag"
	"os"

	"github.com/backend-test-cubi-casa/config"
	"github.com/gin-gonic/gin"
)

// Init system
func init() {
	baseDir := flag.String("e", os.Getenv("HOME_DIR"), "")
	gin.SetMode(gin.TestMode)
	config.Init("development", *baseDir)
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}
