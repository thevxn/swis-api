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
		PostProject)
	g.POST("/restore",
		PostDumpRestore)
}
