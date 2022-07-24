package finance

import (
	"github.com/gin-gonic/gin"
)

// finance CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/", GetAccounts)
	g.GET("/accounts/:owner", GetAccountByOwner)
	g.POST("/restore", PostDumpRestore)
}
