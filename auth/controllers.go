package auth

import (
	"log"
	"net/http"
	"os"

	"swis-api/users"

	"github.com/gin-gonic/gin"
)

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{
		"message": message,
		"code":    code,
	})
}

// https://sosedoff.com/2014/12/21/gin-middleware.html
func AuthMiddleware() gin.HandlerFunc {
	rootToken := os.Getenv("ROOT_TOKEN")

	if rootToken == "" {
		log.Fatal("root token cannot be blank!")
	}

	// flush params --- does nothing
	Params.BearerToken = ""

	return func(c *gin.Context) {
		Params.BearerToken = c.Request.Header.Get("X-Auth-Token")
		//c.ShouldBindHeader(&Params)

		// empty token is disallowed
		if Params.BearerToken == "" {
			respondWithError(c, http.StatusUnauthorized, "empty token")
			return
		}

		// try root token
		if Params.BearerToken == rootToken {
			// pass root name and continue
			Params.User = users.User{Name: "root"}
			c.Next()
			return
		}

		// look for token's non-root owner
		if authUser := users.FindUserByToken(Params.BearerToken); authUser == nil {
			respondWithError(c, http.StatusUnauthorized, "invalid token")
			return
		} else {
			// found, ergo assign that user to auth context
			Params.User = *authUser
		}

		//c.Next()
	}
}
