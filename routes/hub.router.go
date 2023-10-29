package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/backend-test-cubi-casa/api/v1/controllers"
	"github.com/backend-test-cubi-casa/api/v1/repository"
	"github.com/backend-test-cubi-casa/api/v1/services"
)

type HubRoutes struct {
	db *gorm.DB
}

func NewHubRoutes(db *gorm.DB) HubRoutes {
	return HubRoutes{
		db: db,
	}
}

func (rc *HubRoutes) HubRoute(rg *gin.RouterGroup) {
	// Inject dependency
	repo := repository.NewHubRepository(rc.db)
	srv := services.NewHubService(repo)
	ctrl := controllers.NewHubController(srv)

	apiV1 := rg.Group("/v1")
	{
		apiV1.POST("/hubs", ctrl.Create)
	}
}
