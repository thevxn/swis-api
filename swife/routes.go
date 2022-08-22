package swife

import (
	"github.com/gin-gonic/gin"
)

// swife CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/:sitename", GetFrontendBySiteName)
}