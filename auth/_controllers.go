package auth

import (
	"log"
	"net/http"
	"time"

	"swis-api/users"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const TokenExpireDuration = time.Hour * 2

// SetAuthHeaders
func SetAuthHeaders(c *gin.Context) (_ *AuthParams) {
	if err := c.ShouldBindHeader(&Params); err != nil {
		//c.H(500, err)
	}

	return &Params
}

func Authenticate(c *gin.Context) {
	// check for such given token
	if tokens := users.GetUserTokens(); tokens != nil {
		//log.Print(*tokens)
		for _, token := range *tokens {
			if &token != nil && Params.BearerToken == token {
				// ok, do continue
				return
				break
			}
		}
	}

	c.IndentedJSON(http.StatusUnauthorized, gin.H{
		"code":    http.StatusUnauthorized,
		"message": "wrong token given, not permited, logging...",
	})
}

// swapi custom Auth middleware
func SwapiAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//t := time.Now()

		c.Set("auth", "")
		c.Next()

		cCtx := c.Copy()

		// flush and reload Params
		Params = AuthParams{}
		SetAuthHeaders(cCtx)
		Authenticate(cCtx)

		// never executed
		status := c.Writer.Status()
		log.Println(status)
	}
}

// gin jwt

func HelloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(IdentityKey)
	c.JSON(200, gin.H{
		"userID":   claims[IdentityKey],
		"userName": user.(*User).UserName,
		"text":     "Hello World.",
	})
}

func PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*User); ok {
		return jwt.MapClaims{
			IdentityKey: v.UserName,
		}
	}
	return jwt.MapClaims{}
}

func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &User{
		UserName: claims[IdentityKey].(string),
	}
}

func Authenticator(c *gin.Context) (interface{}, error) {
	var loginVals login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	userID := loginVals.Username
	password := loginVals.Password

	if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
		return &User{
			UserName:  userID,
			LastName:  "lmao",
			FirstName: "umom",
		}, nil
	}

	return nil, jwt.ErrFailedAuthentication
}

func Authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*User); ok && v.UserName == "admin" {
		return true
	}

	return false
}

func Unauthorized(c *gin.Context, code int, message string) {
	c.IndentedJSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
