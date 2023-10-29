package controllers

import (
	"net/http"

	"errors"

	"github.com/backend-test-cubi-casa/api/v1/models"
	"github.com/backend-test-cubi-casa/api/v1/services"
	"github.com/backend-test-cubi-casa/helpers/errc"
	"github.com/backend-test-cubi-casa/helpers/resp"
	"github.com/backend-test-cubi-casa/helpers/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		c.JSON(http.StatusBadRequest, resp.BadRequest(errc.GetValidationErrMgs(err)))
		return
	}

	data, err := ctrl.srv.HandleCreate(req)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
			errors.Is(err, gorm.ErrForeignKeyViolated) {
			c.JSON(http.StatusBadRequest, resp.BadRequest(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, resp.InternalServerError())
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
	data, err := ctrl.srv.HandleSearch(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.InternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp.Success(data, http.StatusOK))
}
