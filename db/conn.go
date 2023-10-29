package db

import (
	"fmt"

	"github.com/backend-test-cubi-casa/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Init() (*gorm.DB, error) {
	var err error
	if db == nil {
		cfg := config.GetConfig()
		host := cfg.GetString("database.host")
		port := cfg.GetString("database.port")
		user := cfg.GetString("database.username")
		password := cfg.GetString("database.password")
		database := cfg.GetString("database.db_name")

		psqlInfo := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, database,
		)
		db, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
			//TranslateError: true,
			Logger: logger.Default.LogMode(logger.Info),
		})
	}

	return db, err
}

func GetInstance() *gorm.DB {
	return db
}
