package depots

import (
	"github.com/gin-gonic/gin"
)

// depot CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetAllDepotItems)
	g.POST("/restore",
		PostDumpRestore)

	g.POST("/items/:key",
		PostNewDepotItemByKey)
	g.GET("/items/:key",
		GetDepotItemByKey)
	g.GET("/items/owner/:owner",
		GetDepotItemsByOwner)
	g.PUT("/items/:key",
		UpdateDepotItemByKey)
	g.DELETE("/items/:key",
		DeleteDepotItemByKey)
	g.GET("/types",
		ListTypes)
}
