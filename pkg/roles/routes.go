package roles

import (
	"github.com/gin-gonic/gin"
)

// roles CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetRoles)
	g.POST("/:key",
		PostNewRoleByKey)
	g.GET("/:key",
		GetRoleByKey)
	g.PUT("/:key",
		UpdateRoleByKey)
	g.DELETE("/:key",
		DeleteRoleByKey)
	g.POST("/restore",
		PostDumpRestore)
	g.GET("/types",
		ListTypes)
}
