package finance

import (
	"github.com/gin-gonic/gin"
)

// finance CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetRootData)
	g.POST("/restore",
		PostDumpRestore)

	g.GET("/accounts/",
		GetAccounts)
	g.GET("/accounts/:key",
		GetAccountByKey)
	g.GET("/accounts/owner/:key",
		GetAccountByOwnerKey)
	g.POST("/accounts/:key",
		PostNewAccountByKey)
	g.PUT("/accounts/:key",
		UpdateAccountByKey)
	g.DELETE("/accounts/:key",
		DeleteAccountByKey)

	g.GET("/items/",
		GetItems)
	g.GET("/items/:key",
		GetItemByKey)
	g.GET("/items/account/:key",
		GetItemsByAccountID)
	g.POST("/items/:key",
		PostNewItemByKey)
	g.PUT("/items/:key",
		UpdateItemByKey)
	g.DELETE("/items/:key",
		DeleteItemByKey)

	g.GET("/taxes/:owner/:year",
		DoTaxesByOwner)
}
