package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/backend-test-cubi-casa/config"
	"github.com/backend-test-cubi-casa/db"
	"github.com/backend-test-cubi-casa/helpers/mytime"
	"github.com/backend-test-cubi-casa/middleware"
	"github.com/backend-test-cubi-casa/routes"

	"github.com/gin-gonic/gin"
)

var (
	server     *gin.Engine
	userRoutes routes.UserRoutes
	hubRoutes  routes.HubRoutes
	teamRoutes routes.TeamRoutes
)

func init() {
	server = gin.New()
	server.Use(gin.Logger())

	env := flag.String("e", os.Getenv("APP_ENV"), "")
	flag.Usage = func() {
		log.Println("Usage: server -e {mode}")
		os.Exit(1)
	}

	flag.Parse()
	config.Init(*env, path.Base("config"))

	// DB connection
	dbConn, err := db.Init()
	if err != nil {
		log.Println("Could not connect to postgres db: %v", err)
	}
	log.Println("DB connected!")

	userRoutes = routes.NewUserRoutes(dbConn)
	teamRoutes = routes.NewTeamRoutes(dbConn)
	hubRoutes = routes.NewHubRoutes(dbConn)

	server = gin.Default()
}

func main() {
	cfg := config.GetConfig()
	port := cfg.GetString("server.port")

	// Set timezone server
	_ = mytime.SetTimezone(cfg.GetString("timezone"))

	// Init route
	router := server.Group("")
	router.GET("/health-check", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "welcome to api")
	})
	router.Use(middleware.ValidateHeader()) // Validate header
	router.Use(middleware.Recovery())       // Customize panic log
	router.Use(middleware.RequestLog())     // format log request -response

	userRoutes.UserRoute(router)
	teamRoutes.TeamRoute(router)
	hubRoutes.HubRoute(router)

	// Start server
	_ = server.Run(":" + port)
}
