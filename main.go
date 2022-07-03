// go:build ignore

// Package swis-api is RESTful API core backend for sakalWeb Information System v5.
package main

import (
	"net/http"

	// swapi modules
	"swis-api/alvax"
	"swis-api/auth"
	"swis-api/business"
	"swis-api/depot"
	"swis-api/dish"
	"swis-api/finance"
	//"swis-api/flower"
	"swis-api/groups"
	"swis-api/infra"
	"swis-api/news"
	"swis-api/users"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// reqs from this IPs are treated as proxies, ergo log the real client IP address
	/*swapiProxies := []string{
		"10.4.5.130/25",
	}*/

	//router.SetTrustedProxies(swapiProxies)

	// root path --- testing Bearer print TODO: delete this
	router.GET("/", func(c *gin.Context){
		auth.SetAuthHeaders(c)

		c.JSON(http.StatusOK, gin.H{
			"title": "swAPI v5 RESTful root",
			"code": http.StatusOK,
			"message": "welcome to sakalWeb API (swapi) root",
			"bearer": auth.Params.BearerToken,
		})
	})
	// default 404 route
	router.NoRoute(func(c *gin.Context){
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"message": "unknown route",
		})
	})

	// alvax CRUD
	router.GET("/alvax/commands", alvax.GetCommandList)
	router.POST("/alvax/commands/restore", alvax.PostDumpRestore)

	// business CRUD
	router.GET("/business", business.GetBusinessArray)
	router.POST("/business", business.PostBusiness)
	router.POST("/business/restore", business.PostDumpRestore)

	// depot CRUD
	router.GET("/depots", depot.GetDepots)
	router.GET("/depots/:owner", depot.GetDepotByOwner)
	router.POST("/depots/restore", depot.PostDumpRestore)
	//router.GET("/depots/:groupID", depot.GetDepotByGroupID)
	//router.GET("/depots/:userID", depot.GetDepotByUserID)

	// dish CRUD
	router.HEAD("/dish/test", dish.HeadTest)
	router.GET("/dish/sockets", dish.GetSocketList)
	router.GET("/dish/sockets/:host", dish.GetSocketListByHost)
	router.POST("/dish/sockets/restore", dish.PostDumpRestore)

	// finance CRUD
	router.GET("/finance", finance.GetAccounts)
	router.GET("/finance/accounts/:owner", finance.GetAccountByOwner)
	router.POST("/finance/restore", finance.PostDumpRestore)

	// groups CRUD
	router.GET("/groups", groups.GetGroups)
	router.GET("/groups/:id", groups.GetGroupByID)
	router.POST("/groups", groups.PostGroup)
	//router.PUT("/groups/:id", groups.PutGroupByID)
	//router.DELETE("/groups/:id", groups.DeleteGroupByID)

	// infra CRUD
	router.GET("/infra", infra.GetInfrastructure)
	router.GET("/infra/hosts", infra.GetHosts)
	router.GET("/infra/networks", infra.GetNetworks)
	router.GET("/infra/hosts/:hostname", infra.GetHostByHostname)
	//router.GET("/infra/hosts/:hyp/vms", infra.GetVirtualsByHypervisorName)
	router.POST("/infra/restore", infra.PostDumpRestore)

	// news CRUD
	router.GET("/news/:user", news.GetNewsByUser)
	router.GET("/news/sources", news.GetSources)

	// users CRUD
	router.GET("/users", users.GetUsers)
	router.GET("/users/:name", users.GetUserByName)
	router.POST("/users", users.PostNewUser)
	router.POST("/users/:name/keys/ssh", users.PostUsersSSHKeys)
	router.GET("/users/:name/keys/ssh", users.GetUsersSSHKeysRaw)
	router.POST("/users/restore", users.PostDumpRestore)
	//router.PUT("/users/:id", users.PutUserByID)
	//router.DELETE("/users/:id", users.DeleteUserByID)

	// attach router to http.Server and start it
	router.Run(":8080")
}
