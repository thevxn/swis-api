package business

import (
	"github.com/gin-gonic/gin"
)

// business CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetBusinessEntities)
	g.POST("/",
		PostBusiness)
	g.GET("/:id",
		GetBusinessByID)
	g.PUT("/:id",
		UpdateBusinessByID)
	g.DELETE("/:id",
		DeleteBusinessByID)
	g.POST("/restore",
		PostDumpRestore)
}
