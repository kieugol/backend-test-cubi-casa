package controllers

import (
	"net/http"

	"github.com/backend-test-cubi-casa/api/v1/models"
	"github.com/backend-test-cubi-casa/api/v1/services"
	"github.com/backend-test-cubi-casa/helpers/error"
	errC "github.com/backend-test-cubi-casa/helpers/error"
	"github.com/backend-test-cubi-casa/helpers/resp"
	"github.com/gin-gonic/gin"
)

type HubController struct {
	srv services.IHubService
}

func NewHubController(srv services.IHubService) *HubController {
	return &HubController{
		srv: srv,
	}
}

func (ctrl HubController) Create(c *gin.Context) {
	var req models.HubCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, resp.BadRequest(error.GetValidationErrMgs(err)))
		return
	}

	data, err := ctrl.srv.HandleCreate(c, req)
	if err != nil {
		if errPg := errC.GetGORMErrMgs(err); errPg != nil {
			c.JSON(http.StatusBadRequest, resp.BadRequest(errPg))
			return
		}
		c.JSON(http.StatusInternalServerError, resp.InternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp.Success(data, http.StatusOK))
}
