// @title swis-api v5
// @version 5.4.23
// @description sakalWeb Information System v5 RESTful API documentation
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://savla.dev/swapi
// @contact.email krusty@savla.dev

// @license.name MIT
// @license.url https://github.com/savla-dev/swis-api/blob/master/LICENSE

// @host swis-api-run:8050
// @BasePath /

// @securityDefinitions.apikey apiKey
// @type apiKey
// @name X-Auth-Token
// @in header
// @securityScheme authRequired

// Package swis-api is RESTful API core backend aka 'sakalWeb Information System v5'.
// Basically it is a system of high modularity, where each module (package in golang terminology)
// has its routes, models, and controllers (handler functions) defined in its own directory.
package main

import (
	// golang libs
	"log"
	"net/http"
	"os"
	"time"

	// swapi modules -- very local dependencies
	"go.savla.dev/swis/v5/auth"
	"go.savla.dev/swis/v5/backups"
	"go.savla.dev/swis/v5/business"
	"go.savla.dev/swis/v5/config"
	"go.savla.dev/swis/v5/depots"
	"go.savla.dev/swis/v5/dish"
	"go.savla.dev/swis/v5/finance"
	"go.savla.dev/swis/v5/infra"
	"go.savla.dev/swis/v5/links"
	"go.savla.dev/swis/v5/news"
	"go.savla.dev/swis/v5/projects"
	"go.savla.dev/swis/v5/roles"
	"go.savla.dev/swis/v5/system"
	"go.savla.dev/swis/v5/users"

	// remote dependencies
	gin "github.com/gin-gonic/gin"
)

func main() {
	// blank gin without any middleware
	//gin.DisableConsoleColor()
	router := gin.New()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	//router.Use(gin.Logger())
	router.Use(config.JSONLogMiddleware())

	// CORS Middleware
	router.Use(config.CORSMiddleware())

	// serve savla-dev internal favicon
	router.StaticFile("/favicon.ico", "./favicon.ico")

	// (GET /ping)
	// @Summary Simple ping-pong route
	// @Description Simple ping-pong route
	// @Success 200
	// @Router /ping [get]
	// very simple LE support --- https://github.com/gin-gonic/gin#support-lets-encrypt
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// use custom swapi Auth middlewares --- token auth
	// AuthenticationMiddleware takes care of token verification against loaded Users data structure.
	// AuthorizationMiddleware checks Access Control List (ACL) and the right for dangerous methods usage.
	router.Use(auth.AuthenticationMiddleware())
	router.Use(auth.AuthorizationMiddleware())

	// root path
	router.GET("/", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{
			"header_title": "sakalWebIS v5 RESTful API -- root route",
			"message":      "welcome to swis, " + auth.Params.User.Name + "!",
			"code":         http.StatusOK,
			"app_env": gin.H{
				"app_mode_environment": os.Getenv("APP_ENVIRONMENT"),
				"instance_name":        os.Getenv("HOSTNAME"),
				"alpine_version":       os.Getenv("ALPINE_VERSION"),
				"app_version":          os.Getenv("APP_VERSION"),
				"golang_version":       os.Getenv("GOLANG_VERSION"),
			},
			"timestamp": time.Now().Unix(),
			"user": gin.H{
				"acl":   auth.Params.ACL,
				"roles": auth.Params.Roles,
			},
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

	// system CRUD
	systemRouter := router.Group("/system")
	system.Routes(systemRouter)

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
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		// = 1 * 2^20 = 1,048,576
		MaxHeaderBytes: 1 << 20,
		// use config.CORSMiddleware()
		DisableGeneralOptionsHandler: true,
	}
	server.ListenAndServe()
}
