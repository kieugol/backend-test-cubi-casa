package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/backend-test-cubi-casa/api/v1/controllers"
	"github.com/backend-test-cubi-casa/api/v1/repository"
	"github.com/backend-test-cubi-casa/api/v1/services"
)

type TeamRoutes struct {
	db *gorm.DB
}

func NewTeamRoutes(db *gorm.DB) TeamRoutes {
	return TeamRoutes{
		db: db,
	}
}

func (rc *TeamRoutes) TeamRoute(rg *gin.RouterGroup) {
	// Inject dependency
	repo := repository.NewTeamRepository(rc.db)
	srv := services.NewTeamService(repo)
	ctrl := controllers.NewTeamController(srv)

	apiV1 := rg.Group("/v1")
	{
		apiV1.GET("/teams", ctrl.Search)
		apiV1.POST("/teams", ctrl.Create)
	}
}
