package depots

import (
	"github.com/gin-gonic/gin"
)

// depot CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetDepots)
	g.POST("/",
		PostNewDepot)
	g.GET("/:owner",
		GetDepotByOwner)
	g.DELETE("/:owner",
		DeleteDepotByOwner)
	g.POST("/restore",
		PostDumpRestore)
}
