package alvax

import (
	"github.com/gin-gonic/gin"
)

// alvax CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetConfigs)
	g.GET("/:key",
		GetConfigByKey)
	g.POST("/:key",
		PostNewConfigByKey)
	g.PUT("/:key",
		UpdateConfigByKey)
	g.DELETE("/:key",
		DeleteConfigByKey)
	g.POST("/restore",
		PostDumpRestore)
}
