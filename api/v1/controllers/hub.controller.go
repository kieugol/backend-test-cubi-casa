package controllers

import (
	"net/http"

	"errors"

	"github.com/backend-test-cubi-casa/api/v1/models"
	"github.com/backend-test-cubi-casa/api/v1/services"
	"github.com/backend-test-cubi-casa/helpers/errc"
	"github.com/backend-test-cubi-casa/helpers/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		c.JSON(http.StatusBadRequest, resp.BadRequest(errc.GetValidationErrMgs(err)))
		return
	}

	data, err := ctrl.srv.HandleCreate(c, req)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusBadRequest, resp.BadRequest(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, resp.InternalServerError())
		return
	}

	c.JSON(http.StatusOK, resp.Success(data, http.StatusOK))
}
