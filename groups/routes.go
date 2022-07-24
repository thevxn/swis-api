package groups

import (
	"github.com/gin-gonic/gin"
)

// groups CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/", GetGroups)
	g.GET("/:id", GetGroupByID)
	g.POST("/", PostGroup)
	//g.PUT("/:id", PutGroupByID)
	//g.DELETE("/:id", DeleteGroupByID)
}
