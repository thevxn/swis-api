package users

import (
	"github.com/gin-gonic/gin"
)

// users CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetUsers)
	g.GET("/:key",
		GetUserByKey)
	g.POST("/:key",
		PostNewUserByKey)
	g.PUT("/:key",
		UpdateUserByKey)
	g.DELETE("/:key",
		DeleteUserByKey)
	g.POST("/restore",
		PostDumpRestore)
	g.GET("/types",
		ListTypes)

	g.PUT("/:key/active",
		ActiveToggleUserByKey)
	g.POST("/:key/keys/ssh",
		PostUsersSSHKeys)
	g.GET("/:key/keys/ssh",
		GetUsersSSHKeysRaw)
}
