package news

import (
	"github.com/gin-gonic/gin"
)

// news CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/:user", GetNewsByUser)
	g.GET("/sources", GetSources)
	g.GET("/sources/:user", GetSourcesByUser)
}
