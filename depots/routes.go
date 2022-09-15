package depots

import (
	"github.com/gin-gonic/gin"
)

// depot CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetDepots)
	g.GET("/:owner",
		GetDepotByOwner)
	g.POST("/restore",
		PostDumpRestore)
}
