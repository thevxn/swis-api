// Package swis-api is RESTful API core backend aka 'sakalWeb Information System v5'.
// Basically it is a system of high modularity, where each module (package in golang terminology)
// has its routes, models, and controllers (handler functions) defined in its own folder.
package main

import (
	// golang libs
	"net/http"
	"time"

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

	// remote dependencies
	"github.com/gin-gonic/gin"
)

func main() {
	// blank gin without any middleware
	router := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// by default gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	//authRouter := router.Group("/")

	// root path --- auth required
	router.GET("/", func(c *gin.Context) {
		//auth.SetAuthHeaders(c)

		c.JSON(http.StatusOK, gin.H{
			"title":   "sakalWebIS v5 RESTful API -- root route",
			"message": "welcome to sakalWeb API (swapi) root",
			"code":    http.StatusOK,
			"bearer":  auth.Params.BearerToken,
		})
	})

	// very simple LE support --- https://github.com/gin-gonic/gin#support-lets-encrypt
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// default 404 route
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "unknown route, or disallowed method",
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

	// webui CRUD
	// only a small proposal for simple swapi fronted to browse swapi data quickly
	// import "swis-api/webui"
	//webuiRouter := router.Group("/webui")
	//webui.Routes(webuiRouter)

	// attach router to http.Server and start it
	// https://pkg.go.dev/net/http#Server
	server := &http.Server{
		Addr:         ":8049",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		//MaxHandlerBytes: 1 << 20,
	}
	server.ListenAndServe()
}
