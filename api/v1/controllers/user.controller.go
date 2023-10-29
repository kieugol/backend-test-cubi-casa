package controllers

import (
	"log"
	"net/http"

	"github.com/backend-test-cubi-casa/api/v1/models"
	"github.com/backend-test-cubi-casa/api/v1/services"
	errC "github.com/backend-test-cubi-casa/helpers/error"
	"github.com/backend-test-cubi-casa/helpers/resp"
	"github.com/backend-test-cubi-casa/helpers/util"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	srv services.IUserService
}

func NewUserController(srv services.IUserService) *UserController {
	return &UserController{
		srv: srv,
	}
}

func (ctrl UserController) Create(c *gin.Context) {
	var req models.UserCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, resp.BadRequest(errC.GetValidationErrMgs(err)))
		return
	}

	data, err := ctrl.srv.HandleCreate(c, req)
	if err != nil {
		if errPg := errC.GetGORMErrMgs(err); errPg != nil {
			c.JSON(http.StatusBadRequest, resp.BadRequest(errPg))
			return
		}
		c.JSON(http.StatusOK, resp.Success(data, http.StatusOK))
		return
	}

	c.JSON(http.StatusOK, resp.Success(data, http.StatusOK))
}

func (ctrl UserController) Search(c *gin.Context) {
	var req models.UserSearchReq
	_ = c.ShouldBindQuery(&req)

	if util.IsEmptyStruct(req, models.UserSearchReq{}) {
		c.JSON(http.StatusOK, resp.Success(nil, http.StatusOK))
		return
	}

	log.Println("c ....test:", c)
	data, err := ctrl.srv.HandleSearch(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.InternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp.Success(data, http.StatusOK))
}
