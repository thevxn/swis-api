package alvax

import (
	"github.com/gin-gonic/gin"
)

// alvax CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("",
		GetConfigs)
	g.POST("",
		PostNewConfig)
	g.GET("/:key",
		GetConfigByKey)
	g.PUT("/:key",
		UpdateConfigByKey)
	g.DELETE("/:key",
		DeleteConfigByKey)
	g.POST("/restore",
		PostDumpRestore)
	g.GET("/types",
		ListTypes)
}
