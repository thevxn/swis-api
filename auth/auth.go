package auth

import (
	//"fmt"
	//"net/http"

	"github.com/gin-gonic/gin"
)

type AuthParams struct {
	BearerToken	string `header:"X-Auth-Bearer"`
}

var (
	Params = AuthParams{
		BearerToken: "",
	}
)

//
func SetAuthHeaders(c *gin.Context) (_ *AuthParams) {
	if err := c.ShouldBindHeader(&Params); err != nil {
		c.JSON(500, err)
	}

	return &Params
}
