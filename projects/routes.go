package projects

import (
	"github.com/gin-gonic/gin"
)

// projects CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetProjects)
	g.GET("/:id",
		GetProjectByID)
	g.POST("/",
		PostNewProject)
	g.PUT("/:id",
		UpdateProjectByID)
	g.DELETE("/:id",
		DeleteProjectByID)
	g.POST("/restore",
		PostDumpRestore)
}
