package depots

import (
	"github.com/gin-gonic/gin"
)

// depots CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/items",
		GetAllDepotItems)
	g.GET("/items/types",
		ListTypes)
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

	g.POST("/restore",
		PostDumpRestore)
}
