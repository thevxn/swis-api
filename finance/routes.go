package finance

import (
	"github.com/gin-gonic/gin"
)

// finance CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetAccounts)
	g.POST("/",
		PostNewAccount)
	//g.GET("/accounts/:id",
	//	GetAccountByID)
	g.PUT("/accounts/:idr",
		UpdateAccountByID)
	g.DELETE("/accounts/:id",
		DeleteAccountByID)
	g.GET("/accounts/:owner",
		GetAccountByOwner)
	g.POST("/restore",
		PostDumpRestore)
}
