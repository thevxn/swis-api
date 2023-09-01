package auth

import (
	"log"
	"net/http"
	"os"
	"strings"

	"go.savla.dev/swis/v5/users"

	"github.com/gin-gonic/gin"
)

var Params = AuthParams{
	// Wipe Token string at every request not to allow token forgery.
	BearerToken: "",
}

// https://sosedoff.com/2014/12/21/gin-middleware.html
func AuthenticationMiddleware() gin.HandlerFunc {
	rootToken := os.Getenv("ROOT_TOKEN")

	// stop server if root token environment var is not set
	if rootToken == "" {
		log.Fatal("ROOT_TOKEN environment variable not provided! stopping the server now...")
	}

	return func(c *gin.Context) {
		Params.BearerToken = ""
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

		// look for token's non-root _active_ owner
		if authUser := users.FindUserByToken(Params.BearerToken); authUser == nil {
			respondWithError(c, http.StatusUnauthorized, "invalid token")
			return
		} else {
			// found, ergo assign that user to auth context
			Params.User = *authUser
			Params.Roles = authUser.Roles
			Params.ACL = authUser.ACL
		}

		//c.Next()
	}
}

func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// grant all to root
		if Params.User.Name == "root" {
			c.Next()
			return
		}

		// implement authorization ACL
		for _, item := range Params.ACL {
			path := strings.Split(c.FullPath(), "/")

			// serve root path for everyone
			if c.FullPath() == "/" {
				c.Next()
				return
			}

			// check first requested path "item" against ACL
			if len(path) > 1 && path[1] == item {
				// check the persmission for the requested method usage
				if ok := checkMethodUsagePermission(Params, c); !ok {
					respondWithError(c, http.StatusMethodNotAllowed, "method not allowed")
					return
				}

				// access granted according to ACL item and role by method type
				c.Next()
				return
			}
		}

		respondWithError(c, http.StatusForbidden, "access denied")
		return
	}
}
