package users

import (
	"github.com/gin-gonic/gin"
)

// users CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetUsers)
	g.POST("/",
		PostNewUser)
	g.GET("/:name",
		GetUserByName)
	g.PUT("/:name",
		UpdateUserByName)
	g.DELETE("/:name",
		DeleteUserByName)
	g.PUT("/:name/active",
		ActiveToggleUserByName)
	g.POST("/:name/keys/ssh",
		PostUsersSSHKeys)
	g.GET("/:name/keys/ssh",
		GetUsersSSHKeysRaw)
	g.POST("/restore",
		PostDumpRestore)
}
