package business

import (
	"github.com/gin-gonic/gin"
)

// business CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetBusinessEntities)
	g.POST("/:key",
		PostBusinessByKey)
	g.GET("/:key",
		GetBusinessByKey)
	g.PUT("/:key",
		UpdateBusinessByKey)
	g.DELETE("/:key",
		DeleteBusinessByKey)
	g.POST("/restore",
		PostDumpRestore)
	g.GET("/types",
		ListTypes)
}
