package unit

import (
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

func TestCreateTeam_Success(t *testing.T) {
	// Mock data
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	timeNow := time.Now()
	teamReq := models.TeamCreateReq{
		Name:  "Hub1",
		HubID: 2,
	}
	teamResp := models.Team{
		ID:        1,
		Name:      teamReq.Name,
		HubID:     teamReq.HubID,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	team := sqlmock.
		NewRows([]string{"id", "name", "hub_id", "created_at", "updated_at"}).
		AddRow(teamResp.ID, teamReq.Name, teamReq.HubID, timeNow, timeNow)
	expectedSQLUser := "INSERT INTO \"teams\" (.+) VALUES (.+)"
	mock.ExpectBegin()
	mock.ExpectQuery(expectedSQLUser).WillReturnRows(team)
	mock.ExpectCommit()

	// execute test
	repo := repository.NewTeamRepository(db)
	srv := services.NewTeamService(repo)
	result, err := srv.HandleCreate(teamReq)
	// Assert test result
	assert.Equal(t, nil, err)
	assert.Equal(t, teamResp, *result)
}

func TestCreateTeam_Failed_ForeignKey_NotFound(t *testing.T) {
	// Mock data
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	teamReq := models.TeamCreateReq{
		Name:  "hub-duplicate",
		HubID: 10,
	}
	var teamResp *models.Team

	expectedSQLUser := "INSERT INTO \"teams\" (.+) VALUES (.+)"
	mock.ExpectBegin()
	mock.ExpectQuery(expectedSQLUser).WillReturnError(gorm.ErrForeignKeyViolated)
	mock.ExpectRollback()

	// execute test
	repo := repository.NewTeamRepository(db)
	srv := services.NewTeamService(repo)
	result, err := srv.HandleCreate(teamReq)

	// Assert test result
	assert.Equal(t, gorm.ErrForeignKeyViolated, err)
	assert.Equal(t, teamResp, result)
}

func TestSearchTeam_Found_WithNameAndID(t *testing.T) {
	// Mock data
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	teamReq := models.TeamSearchReq{
		Name: "php",
		ID:   1,
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

	team := sqlmock.
		NewRows([]string{"id", "name", "hub_id"}).
		AddRow(teamResp.ID, teamResp.Name, teamResp.HubID)
	hub := sqlmock.
		NewRows([]string{"id", "name", "location"}).
		AddRow(hubResp.ID, hubResp.Name, hubResp.Location)
	expectedSQLTeam := "SELECT (.+) FROM \"teams\" WHERE id = (.+) AND name LIKE (.+)"
	expectedSQLHub := "SELECT (.+) FROM \"hubs\" WHERE \"hubs\".\"id\" = (.+)"
	mock.ExpectQuery(expectedSQLTeam).WillReturnRows(team)
	mock.ExpectQuery(expectedSQLHub).WillReturnRows(hub)

	// execute test
	repo := repository.NewTeamRepository(db)
	srv := services.NewTeamService(repo)
	result, err := srv.HandleSearch(teamReq)

	// Assert test result
	assert.Equal(t, nil, err)
	assert.Equal(t, teamResp, *result[0])
}

func TestSearchUser_InternalServerError(t *testing.T) {
	// Mock data
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	teamReq := models.TeamSearchReq{
		Name: "team-not-found",
	}
	var teamResp []*models.TeamResp

	expectedSQLUser := "SELECT (.+) FROM \"teams\" WHERE name LIKE (.+)"
	mock.ExpectQuery(expectedSQLUser).WillReturnError(gorm.ErrUnsupportedRelation)

	// execute test
	repo := repository.NewTeamRepository(db)
	srv := services.NewTeamService(repo)
	result, err := srv.HandleSearch(teamReq)

	// Assert test result
	assert.Equal(t, gorm.ErrUnsupportedRelation, err)
	assert.Equal(t, teamResp, result)
}
