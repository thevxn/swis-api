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
	alvaxRouter := router.Group("/alvax")
	alvax.Routes(alvaxRouter)

	// business CRUD
	businessRouter := router.Group("/business")
	business.Routes(businessRouter)

	// depot/depots CRUD
	depotRouter := router.Group("/depots")
	depot.Routes(depotRouter)

	// dish CRUD
	dishRouter := router.Group("/dish")
	dish.Routes(dishRouter)

	// finance accounts CRUD
	financeRouter := router.Group("/finance")
	finance.Routes(financeRouter)

	// groups CRUD
	groupsRouter := router.Group("/groups")
	groups.Routes(groupsRouter)

	// infra CRUD
	infraRouter := router.Group("/infra")
	infra.Routes(infraRouter)

	// news CRUD
	newsRouter := router.Group("/news")
	news.Routes(newsRouter)

	// projects CRUD
	projectsRouter := router.Group("/projects")
	projects.Routes(projectsRouter)

	// users CRUD
	usersRouter := router.Group("/users")
	users.Routes(usersRouter)

	// attach router to http.Server and start it
	router.Run(":8080")
}
