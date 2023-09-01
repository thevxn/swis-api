package depots

import (
	"github.com/gin-gonic/gin"
)

// depot CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetDepotItems)
	g.POST("/:key",
		PostNewDepotItemByKey)
	g.GET("/:owner",
		GetDepotItemsByOwner)
	g.PUT("/:key",
		UpdateDepotItemByKey)
	g.DELETE("/:key",
		DeleteDepotItemByKey)
	g.POST("/restore",
		PostDumpRestore)
}
