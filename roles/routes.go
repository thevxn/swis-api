package roles

import (
	"github.com/gin-gonic/gin"
)

// roles CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/", GetRoles)
	g.GET("/:id", GetRoleByName)
	g.POST("/", PostRole)
	//g.PUT("/:id", PutRoleByName)
	//g.DELETE("/:id", DeleteRoleByName)
	g.POST("/restore", PostDumpRestore)
}
