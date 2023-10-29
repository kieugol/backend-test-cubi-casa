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

func TestCreateUser_Success(t *testing.T) {
	// Mock data
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	timeNow := time.Now()
	userReq := models.UserCreateReq{
		Name:    "user",
		Email:   "user@gmail.com",
		Phone:   "+84917474772",
		Address: "HCM city",
		TeamID:  1,
	}

	userResp := models.User{
		ID:        1,
		Name:      userReq.Name,
		Email:     userReq.Email,
		Phone:     userReq.Phone,
		Address:   userReq.Address,
		TeamID:    userReq.TeamID,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	user := sqlmock.
		NewRows([]string{"id", "name", "email", "phone", "address", "team_id", "created_at", "updated_at"}).
		AddRow(userResp.ID, userReq.Name, userReq.Email, userReq.Phone, userReq.Address, userReq.TeamID, timeNow, timeNow)
	expectedSQLUser := "INSERT INTO \"users\" (.+) VALUES (.+)"

	mock.ExpectBegin()
	mock.ExpectQuery(expectedSQLUser).WillReturnRows(user)
	mock.ExpectCommit()

	// execute test
	repo := repository.NewUserRepository(db)
	srv := services.NewUserService(repo)
	result, err := srv.HandleCreate(context.TODO(), userReq)

	// Assert test result
	assert.Equal(t, nil, err)
	assert.Equal(t, userResp, *result)
}

func TestCreateUser_Failed_Duplicate(t *testing.T) {
	// Mock data
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	userReq := models.UserCreateReq{
		Name:    "user",
		Email:   "user@gmail.com",
		Phone:   "+84917474772",
		Address: "HCM city",
		TeamID:  1,
	}
	var userResp *models.User

	expectedSQLUser := "INSERT INTO \"users\" (.+) VALUES (.+)"

	mock.ExpectBegin()
	mock.ExpectQuery(expectedSQLUser).WillReturnError(gorm.ErrDuplicatedKey)
	mock.ExpectRollback()

	// execute test
	repo := repository.NewUserRepository(db)
	srv := services.NewUserService(repo)
	result, err := srv.HandleCreate(context.TODO(), userReq)

	// Assert test result
	assert.Equal(t, gorm.ErrDuplicatedKey, err)
	assert.Equal(t, userResp, result)
}

func TestSearchUser_Found_WithNameAndEmail(t *testing.T) {
	// Mock data
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	userReq := models.UserSearchReq{
		Name:  "user",
		Email: "user.acc",
	}
	hubResp := models.HubResp{
		ID:       1,
		Name:     "hub-1",
		Location: "hub-hcm",
	}
	teamResp := models.TeamResp{
		ID:    1,
		Name:  "team-php",
		HubID: 1,
		Hub:   hubResp,
	}
	userResp := models.UserResp{
		ID:      1,
		Name:    "user",
		Email:   "user.acc@gmail.com",
		Phone:   "+84917474772",
		Address: "HCM city",
		TeamID:  1,
		Team:    teamResp,
	}

	user := sqlmock.
		NewRows([]string{"id", "name", "email", "phone", "address", "team_id"}).
		AddRow(userResp.ID, userResp.Name, userResp.Email, userResp.Phone, userResp.Address, userResp.TeamID)
	team := sqlmock.
		NewRows([]string{"id", "name", "hub_id"}).
		AddRow(teamResp.ID, teamResp.Name, teamResp.HubID)
	hub := sqlmock.
		NewRows([]string{"id", "name", "location"}).
		AddRow(hubResp.ID, hubResp.Name, hubResp.Location)
	expectedSQLUser := "SELECT (.+) FROM \"users\" WHERE name LIKE (.+) AND email LIKE (.+)"
	expectedSQLTeam := "SELECT (.+) FROM \"teams\" WHERE \"teams\".\"id\" = (.+)"
	expectedSQLHub := "SELECT (.+) FROM \"hubs\" WHERE \"hubs\".\"id\" = (.+)"
	mock.ExpectQuery(expectedSQLUser).WillReturnRows(user)
	mock.ExpectQuery(expectedSQLTeam).WillReturnRows(team)
	mock.ExpectQuery(expectedSQLHub).WillReturnRows(hub)

	// execute test
	repo := repository.NewUserRepository(db)
	srv := services.NewUserService(repo)
	result, err := srv.HandleSearch(context.TODO(), userReq)

	// Assert test result
	assert.Equal(t, nil, err)
	assert.Equal(t, userResp, *result[0])
}

func TestSearchUser_InternalServerError_(t *testing.T) {
	// Mock data
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	userReq := models.UserSearchReq{
		Name: "user-not-found",
	}
	var userResp []*models.UserResp

	expectedSQLUser := "SELECT (.+) FROM \"users\" WHERE name LIKE (.+)"
	mock.ExpectQuery(expectedSQLUser).WillReturnError(gorm.ErrUnsupportedRelation)

	// execute test
	repo := repository.NewUserRepository(db)
	srv := services.NewUserService(repo)
	result, err := srv.HandleSearch(context.TODO(), userReq)

	// Assert test result
	assert.Equal(t, gorm.ErrUnsupportedRelation, err)
	assert.Equal(t, userResp, result)
}
