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

type TeamController struct {
	srv services.ITeamService
}

func NewTeamController(srv services.ITeamService) *TeamController {
	return &TeamController{
		srv: srv,
	}
}

func (ctrl TeamController) Create(c *gin.Context) {
	var req models.TeamCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, resp.BadRequest(errc.GetValidationErrMgs(err)))
		return
	}

	data, err := ctrl.srv.HandleCreate(c, req)
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

func (ctrl TeamController) Search(c *gin.Context) {
	var req models.TeamSearchReq
	_ = c.ShouldBindQuery(&req)

	if util.IsEmptyStruct(req, models.TeamSearchReq{}) {
		c.JSON(http.StatusOK, resp.Success(nil, http.StatusOK))
		return
	}

	data, err := ctrl.srv.HandleSearch(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.InternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp.Success(data, http.StatusOK))
}
