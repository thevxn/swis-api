package auth

import (
	//"fmt"
	//"net/http"

	"github.com/gin-gonic/gin"
)

// SetAuthHeaders
func SetAuthHeaders(c *gin.Context) (_ *AuthParams) {
	if err := c.ShouldBindHeader(&Params); err != nil {
		c.JSON(500, err)
	}

	return &Params
}
