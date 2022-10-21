// @title swis-api v5
// @version 5.2.15
// @description sakalWeb Information System v5 RESTful API documentation
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://savla.dev/swapi
// @contact.email krusty@savla.dev

// @license.name MIT
// @license.url https://github.com/savla-dev/swis-api/blob/master/LICENSE

// @host swapi.savla.su:8049
// @BasePath /

// @securityDefinitions.apikey apiKey
// @type apiKey
// @name X-Auth-Token
// @in header
// @securityScheme authRequired

// Package swis-api is RESTful API core backend aka 'sakalWeb Information System v5'.
// Basically it is a system of high modularity, where each module (package in golang terminology)
// has its routes, models, and controllers (handler functions) defined in its own folder.
package main

import (
	// golang libs
	"log"
	"net/http"
	"os"
	"time"

	// swapi modules -- very local dependencies
	"swis-api/alvax"
	"swis-api/auth"
	"swis-api/backups"
	"swis-api/business"
	"swis-api/depots"
	"swis-api/dish"
	"swis-api/finance"
	"swis-api/infra"
	"swis-api/links"
	"swis-api/news"
	"swis-api/projects"
	"swis-api/roles"
	"swis-api/six"
	"swis-api/swife"
	"swis-api/users"

	// remote dependencies
	gin "github.com/gin-gonic/gin"
)

func main() {
	// blank gin without any middleware
	gin.DisableConsoleColor()
	router := gin.New()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	router.Use(gin.Logger())

	// CORS Middleware
	router.Use(CORSMiddleware())

	// serve savla-dev internal favicon
	router.StaticFile("/favicon.ico", "./.assets/favicon.ico")

	// (GET /ping)
	// @Summary Simple ping-pong route
	// @Description Simple ping-pong route
	// @Success 200
	// @Router /ping [get]
	// very simple LE support --- https://github.com/gin-gonic/gin#support-lets-encrypt
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// use custom swapi Auth middleware --- token auth
	router.Use(auth.AuthMiddleware())

	// root path
	router.GET("/", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{
			"title":       "sakalWebIS v5 RESTful API -- root route",
			"message":     "welcome to swis, " + auth.Params.User.Name + "!",
			"code":        http.StatusOK,
			"environment": os.Getenv("APP_ENVIRONMENT"),
			"version":     os.Getenv("APP_VERSION"),
			"timestamp":   time.Now().Unix(),
		})
	})

	// default 404 route
	router.NoRoute(func(c *gin.Context) {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "unknown route, or disallowed method",
		})
	})

	//
	// swis modules
	//

	// alvax CRUD
	alvaxRouter := router.Group("/alvax")
	alvax.Routes(alvaxRouter)

	// backups CRUD
	backupsRouter := router.Group("/backups")
	backups.Routes(backupsRouter)

	// business CRUD
	businessRouter := router.Group("/business")
	business.Routes(businessRouter)

	// depots CRUD
	depotsRouter := router.Group("/depots")
	depots.Routes(depotsRouter)

	// dish CRUD
	dishRouter := router.Group("/dish")
	dish.Routes(dishRouter)

	// finance accounts CRUD
	financeRouter := router.Group("/finance")
	finance.Routes(financeRouter)

	// infra CRUD
	infraRouter := router.Group("/infra")
	infra.Routes(infraRouter)

	// links CRUD
	linksRouter := router.Group("/links")
	links.Routes(linksRouter)

	// news CRUD
	newsRouter := router.Group("/news")
	news.Routes(newsRouter)

	// projects CRUD
	projectsRouter := router.Group("/projects")
	projects.Routes(projectsRouter)

	// roles CRUD
	rolesRouter := router.Group("/roles")
	roles.Routes(rolesRouter)

	// six CRUD
	sixRouter := router.Group("/six")
	six.Routes(sixRouter)

	// swife CRUD
	swifeRouter := router.Group("/swife")
	swife.Routes(swifeRouter)

	// users CRUD
	usersRouter := router.Group("/users")
	//usersRouter.Use(authMiddleware.MiddlewareFunc())
	users.Routes(usersRouter)

	// attach router to http.Server and start it, check for DOCKER_INTERNAL_PORT constant
	if os.Getenv("DOCKER_INTERNAL_PORT") == "" {
		log.Fatal("DOCKER_INTERNAL_PORT environment variable not provided! stopping the server now...")
	}

	// https://pkg.go.dev/net/http#Server
	server := &http.Server{
		Addr:         "0.0.0.0:" + os.Getenv("DOCKER_INTERNAL_PORT"),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		//MaxHandlerBytes: 1 << 20,
	}
	server.ListenAndServe()
}
