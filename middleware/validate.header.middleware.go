package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/backend-test-cubi-casa/helpers/resp"
	"github.com/backend-test-cubi-casa/helpers/util"
)

func ValidateHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !util.ShouldBindHeader(c) {
			c.JSON(http.StatusBadRequest, resp.MissingHeader())
			c.Abort()
			return
		}

		c.Next()
	}
}
