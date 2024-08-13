package projects

import (
	"github.com/gin-gonic/gin"
)

// projects CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("",
		GetProjects)
	g.POST("",
		PostNewProject)
	g.GET("/:key",
		GetProjectByKey)
	g.PUT("/:key",
		UpdateProjectByKey)
	g.DELETE("/:key",
		DeleteProjectByKey)
	g.POST("/restore",
		PostDumpRestore)
	g.GET("/types",
		ListTypes)
}
