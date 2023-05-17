package roles

import (
	"github.com/gin-gonic/gin"
)

// roles CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetRoles)
	g.POST("/",
		PostRole)
	g.GET("/:name",
		GetRoleByName)
	g.PUT("/:name",
		UpdateRoleByName)
	g.DELETE("/:name",
		DeleteRoleByName)
	g.POST("/restore",
		PostDumpRestore)
}
