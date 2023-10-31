package unit

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/backend-test-cubi-casa/api/v1/models"
	"github.com/backend-test-cubi-casa/api/v1/repository"
	"github.com/backend-test-cubi-casa/api/v1/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DbMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		t.Fatal(err)
	}

	return sqldb, gormdb, mock
}

func TestSaveHub_Success(t *testing.T) {
	// Mock data
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	timeNow := time.Now()
	hubReq := models.HubCreateReq{
		Name:     "Hub1",
		Location: "HCM city",
	}

	dataExpected := models.Hub{
		ID:        1,
		Name:      hubReq.Name,
		Location:  hubReq.Location,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	user := sqlmock.
		NewRows([]string{"id", "name", "location", "created_at", "updated_at"}).
		AddRow(dataExpected.ID, hubReq.Name, hubReq.Location, timeNow, timeNow)
	expectedSQLUser := "INSERT INTO \"hubs\" (.+) VALUES (.+)"

	mock.ExpectBegin()
	mock.ExpectQuery(expectedSQLUser).WillReturnRows(user)
	mock.ExpectCommit()

	// execute test
	repo := repository.NewHubRepository(db)
	srv := services.NewHubService(repo)
	result, err := srv.HandleCreate(context.TODO(), hubReq)

	// Assert test result
	assert.Equal(t, nil, err)
	assert.Equal(t, dataExpected, *result)
}

func TestSaveHub_Failed_Duplicate(t *testing.T) {
	// Mock data
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	hubReq := models.HubCreateReq{
		Name:     "hub-duplicate",
		Location: "HCM city",
	}
	var dataExpected *models.Hub

	expectedSQLUser := "INSERT INTO \"hubs\" (.+) VALUES (.+)"

	mock.ExpectBegin()
	mock.ExpectQuery(expectedSQLUser).WillReturnError(gorm.ErrDuplicatedKey)
	mock.ExpectRollback()

	// execute test
	repo := repository.NewHubRepository(db)
	srv := services.NewHubService(repo)
	result, err := srv.HandleCreate(context.TODO(), hubReq)

	// Assert test result
	assert.Equal(t, gorm.ErrDuplicatedKey, err)
	assert.Equal(t, dataExpected, result)
}
