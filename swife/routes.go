package swife

import (
	"github.com/gin-gonic/gin"
)

// swife CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/", GetFrontends)
	g.GET("/:sitename", GetFrontendBySiteName)
	g.POST("/restore", PostDumpRestore)
}
