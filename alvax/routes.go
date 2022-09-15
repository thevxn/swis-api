package alvax

import (
	"github.com/gin-gonic/gin"
)

// alvax CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/commands",
		GetCommandList)
	g.POST("/commands/restore",
		PostDumpRestore)
}
