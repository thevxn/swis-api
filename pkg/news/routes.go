package news

import (
	"github.com/gin-gonic/gin"
)

// news CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/:key",
		GetNewsByUserKey)
	g.GET("/sources",
		GetSources)
	g.GET("/sources/:key",
		GetSourcesByUserKey)
	g.POST("/sources/:key",
		PostNewSourcesByUserKey)
	g.PUT("/sources/:key",
		UpdateSourcesByUserKey)
	g.DELETE("/sources/:key",
		DeleteSourcesByUserKey)
	g.POST("/sources/restore",
		PostDumpRestore)
	g.GET("/types",
		ListTypes)
}
