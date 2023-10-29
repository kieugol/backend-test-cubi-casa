package controllers

import (
	"net/http"

	"github.com/backend-test-cubi-casa/api/v1/models"
	"github.com/backend-test-cubi-casa/api/v1/services"
	errC "github.com/backend-test-cubi-casa/helpers/error"
	"github.com/backend-test-cubi-casa/helpers/resp"
	"github.com/backend-test-cubi-casa/helpers/util"
	"github.com/gin-gonic/gin"
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
