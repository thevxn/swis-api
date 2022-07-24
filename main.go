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
	"swis-api/projects"
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
	router.GET("/", func(c *gin.Context) {
		auth.SetAuthHeaders(c)

		c.JSON(http.StatusOK, gin.H{
			"title":   "swAPI v5 RESTful root",
			"code":    http.StatusOK,
			"message": "welcome to sakalWeb API (swapi) root",
			"bearer":  auth.Params.BearerToken,
		})
	})

	// default 404 route
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "unknown route",
		})
	})

	// serve savla-dev internal favicon
	router.StaticFile("/favicon.ico", "./.assets/favicon.ico")

	// alvax CRUD
	aa := router.Group("/alvax")
	{
		aa.GET("/commands", alvax.GetCommandList)
		aa.POST("/commands/restore", alvax.PostDumpRestore)
	}

	// business CRUD
	biz := router.Group("/business")
	{
		biz.GET("/", business.GetBusinessArray)
		biz.GET("/:ico_id", business.GetBusinessByICO)
		biz.POST("/", business.PostBusiness)
		biz.POST("/restore", business.PostDumpRestore)
	}

	// depot CRUD
	de := router.Group("/depots")
	{
		de.GET("/", depot.GetDepots)
		de.GET("/:owner", depot.GetDepotByOwner)
		de.POST("/restore", depot.PostDumpRestore)
		//de.GET("/:groupID", depot.GetDepotByGroupID)
		//gg.GET("/:userID", depot.GetDepotByUserID)
	}

	// dish CRUD
	di := router.Group("/dish")
	{
		di.HEAD("/test", dish.HeadTest)
		di.GET("/sockets", dish.GetSocketList)
		di.GET("/sockets/:host", dish.GetSocketListByHost)
		//di.POST("/sockets/result", dish.PostSocketTestResult)
		di.POST("/sockets/restore", dish.PostDumpRestore)
	}

	// finance accounts CRUD
	ff := router.Group("/finance")
	{
		ff.GET("/", finance.GetAccounts)
		ff.GET("/accounts/:owner", finance.GetAccountByOwner)
		ff.POST("/restore", finance.PostDumpRestore)
	}

	// groups CRUD
	gg := router.Group("/groups")
	{
		gg.GET("/", groups.GetGroups)
		gg.GET("/:id", groups.GetGroupByID)
		gg.POST("/", groups.PostGroup)
		//gg.PUT("/:id", groups.PutGroupByID)
		//gg.DELETE("/:id", groups.DeleteGroupByID)
	}

	// infra CRUD
	ii := router.Group("/infra")
	{
		ii.GET("/", infra.GetInfrastructure)
		ii.GET("/hosts", infra.GetHosts)
		ii.GET("/hosts/:hostname", infra.GetHostByHostname)
		//ii.GET("/hosts/ansible/:ansible_group", infra.GetHostsByAnsibleGroup)
		//ii.GET("/map", infra.GetInfraMap)
		ii.GET("/networks", infra.GetNetworks)
		//ii.GET("/hosts/:hyp/vms", infra.GetVirtualsByHypervisorName)
		ii.POST("/restore", infra.PostDumpRestore)
	}

	// news CRUD
	nn := router.Group("/news")
	{
		nn.GET("/:user", news.GetNewsByUser)
		nn.GET("/sources", news.GetSources)
	}

	// projects CRUD
	pp := router.Group("/projects")
	{
		pp.GET("/", projects.GetProjects)
		pp.GET("/:id", projects.GetProjectByID)
		pp.POST("/", projects.PostProject)
		pp.POST("/restore", projects.PostDumpRestore)
	}

	// users CRUD
	uu := router.Group("/users")
	{
		uu.GET("/", users.GetUsers)
		uu.GET("/:name", users.GetUserByName)
		uu.POST("/", users.PostNewUser)
		uu.POST("/:name/keys/ssh", users.PostUsersSSHKeys)
		uu.GET("/:name/keys/ssh", users.GetUsersSSHKeysRaw)
		uu.POST("/restore", users.PostDumpRestore)
		//uu.PUT("/:id", users.PutUserByID)
		//uu.DELETE("/:id", users.DeleteUserByID)
	}

	// attach router to http.Server and start it
	router.Run(":8080")
}
