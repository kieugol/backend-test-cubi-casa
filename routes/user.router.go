package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/backend-test-cubi-casa/api/v1/controllers"
	"github.com/backend-test-cubi-casa/api/v1/repository"
	"github.com/backend-test-cubi-casa/api/v1/services"
)

type UserRoutes struct {
	db *gorm.DB
}

func NewUserRoutes(db *gorm.DB) UserRoutes {
	return UserRoutes{
		db: db,
	}
}

func (rc *UserRoutes) UserRoute(rg *gin.RouterGroup) {
	// Inject dependency
	repo := repository.NewUserRepository(rc.db)
	srv := services.NewUserService(repo)
	ctrl := controllers.NewUserController(srv)

	apiV1 := rg.Group("/v1")
	{
		apiV1.GET("/users", ctrl.Search)
		apiV1.POST("/users", ctrl.Create)
	}
}
