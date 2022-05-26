// go:build ignore

// Package swis-api is RESTful API core backend for sakalWeb Information System v5.
package main

import (
	"net/http"

	// swapi modules
	
	"swis-api/alvax"
	"swis-api/auth"
	//"swis-api/b2b"
	"swis-api/depot"
	"swis-api/dish"
	//"swis-api/flower"
	"swis-api/groups"
	//"swis-api/hosts"
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

	// depot CRUD
	router.GET("/depots", depot.GetDepots)
	router.GET("/depots/:owner", depot.GetDepotByOwner)
	router.POST("/depots/restore", depot.PostDepotDumpRestore)
	//router.GET("/depots/:groupID", depot.GetDepotByGroupID)
	//router.GET("/depots/:userID", depot.GetDepotByUserID)

	// dish CRUD
	router.HEAD("/dish/test", dish.HeadTest)
	router.GET("/dish/sockets", dish.GetSocketList)
	router.GET("/dish/sockets/:host", dish.GetSocketListByHost)

	// dump CRUD -- backuping routes -- GET to dump arrays to JSON files, POST to rebuild array from JSON dumps
	//router.GET("/dump/users", users.GetDump)
	//router.POST("/dump/users", users.PostDump)

	// groups CRUD
	router.GET("/groups", groups.GetGroups)
	router.GET("/groups/:id", groups.GetGroupByID)
	router.POST("/groups", groups.PostGroup)
	//router.PUT("/groups/:id", groups.PutGroupByID)
	//router.DELETE("/groups/:id", groups.DeleteGroupByID)

	// hosts CRUD
	//router.GET("/hosts", hosts.GetHosts)
	//router.GET("/hosts/:hostname", hosts.GetHostByHostname)
	//router.GET("/hosts/:hyp/vms", hosts.GetVirtualsByHypervisorName)
	//router.POST("/hosts/restore", hosts.PostHostsDumpRestore)

	// users CRUD
	router.GET("/users", users.GetUsers)
	router.GET("/users/:id", users.GetUserByID)
	router.POST("/users", users.PostUser)
	router.POST("/users/restore", users.PostUsersDumpRestore)
	router.POST("/users/:id/keys/ssh", users.PostUserSSHKey)
	//router.PUT("/users/:id", users.PutUserByID)
	//router.DELETE("/users/:id", users.DeleteUserByID)

	// attach router to http.Server and start it
	router.Run(":8080")
}

