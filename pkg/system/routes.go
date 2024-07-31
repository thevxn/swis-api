package system

import (
	"github.com/gin-gonic/gin"
)

// system CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/packages",
		GetMountedPackages)
}
