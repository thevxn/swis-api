package projects

import (
	"github.com/gin-gonic/gin"
)

// projects CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetProjects)
	g.GET("/:key",
		GetProjectByKey)
	g.POST("/:key",
		PostNewProjectByKey)
	g.PUT("/:key",
		UpdateProjectByKey)
	g.DELETE("/:key",
		DeleteProjectByKey)
	g.POST("/restore",
		PostDumpRestore)
}
