package business

import (
	"github.com/gin-gonic/gin"
)

// business CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/", GetBusinessArray)
	g.GET("/:ico_id", GetBusinessByICO)
	g.POST("/", PostBusiness)
	g.POST("/restore", PostDumpRestore)
}
