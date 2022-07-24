package users

import (
	"github.com/gin-gonic/gin"
)

// users CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/", GetUsers)
	g.GET("/:name", GetUserByName)
	g.POST("/", PostNewUser)
	g.POST("/:name/keys/ssh", PostUsersSSHKeys)
	g.GET("/:name/keys/ssh", GetUsersSSHKeysRaw)
	g.POST("/restore", PostDumpRestore)
	//g.PUT("/:id", PutUserByID)
	//g.DELETE("/:id", DeleteUserByID)
}
