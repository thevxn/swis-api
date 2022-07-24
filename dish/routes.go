package dish

import (
	"github.com/gin-gonic/gin"
)

// dish CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.HEAD("/test", HeadTest)
	g.GET("/sockets", GetSocketList)
	g.GET("/sockets/:host", GetSocketListByHost)
	//g.POST("/sockets/result", PostSocketTestResult)
	g.POST("/sockets/restore", PostDumpRestore)
}
