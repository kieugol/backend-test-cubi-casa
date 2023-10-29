package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/backend-test-cubi-casa/api/v1/models"
	"github.com/backend-test-cubi-casa/config"
	"github.com/backend-test-cubi-casa/db"
	"gorm.io/gorm"
)

var (
	pgdb *gorm.DB
	err  error
)

func init() {
	env := flag.String("e", os.Getenv("APP_ENV"), "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}

	flag.Parse()
	config.Init(*env, path.Base("config"))

	pgdb, err = db.Init()
	if err != nil {
		log.Fatalf("could not connect to postgres db: %v", err)
	}
	log.Println("DB connected!")

}

func main() {
	pgdb.Migrator().DropTable(&models.User{}, &models.Team{}, models.Hub{})
	pgdb.AutoMigrate(&models.User{}, &models.Team{}, models.Hub{})
	fmt.Println("👍 Migration complete")
}
