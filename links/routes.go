package links

import (
	"github.com/gin-gonic/gin"
)

// links CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetLinks)
	g.POST("/",
		PostNewLink)
	g.GET("/:hash",
		GetLinkByHash)
	g.PUT("/:hash",
		UpdateLinkByHash)
	g.DELETE("/:hash",
		DeleteLinkByHash)
	g.PUT("/:hash/active",
		ActiveToggleByHash)
	g.POST("/restore",
		PostDumpRestore)
}
