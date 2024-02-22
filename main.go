// @title swis-api (swapi) v5
// @version 5.13.18
// @description sakalWeb Information System v5 RESTful API documentation
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://savla.dev/swapi
// @contact.email krusty@savla.dev

// @license.name MIT
// @license.url https://github.com/savla-dev/swis-api/blob/master/LICENSE

// @host swis-api-run-prod:8050
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
	"go.savla.dev/swis/v5/pkg/alvax"
	"go.savla.dev/swis/v5/pkg/auth"
	"go.savla.dev/swis/v5/pkg/backups"
	"go.savla.dev/swis/v5/pkg/business"
	"go.savla.dev/swis/v5/pkg/config"
	"go.savla.dev/swis/v5/pkg/core"
	"go.savla.dev/swis/v5/pkg/depots"
	"go.savla.dev/swis/v5/pkg/dish"
	"go.savla.dev/swis/v5/pkg/finance"
	"go.savla.dev/swis/v5/pkg/infra"
	"go.savla.dev/swis/v5/pkg/links"
	"go.savla.dev/swis/v5/pkg/news"
	"go.savla.dev/swis/v5/pkg/projects"
	"go.savla.dev/swis/v5/pkg/roles"
	"go.savla.dev/swis/v5/pkg/users"

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

	// Mirroring Middleware
	//router.Use(core.MirrorMiddleware())

	// serve savla-dev internal favicon
	router.StaticFile("/favicon.ico", "./favicon.ico")

	// @Summary Simple ping-pong route
	// @Description Simple ping-pong route
	// @Success 200
	// @Router /ping [get]
	// Very simple LE support --- https://github.com/gin-gonic/gin#support-lets-encrypt.
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.HEAD("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// default 404 route
	router.NoRoute(func(c *gin.Context) {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "unknown route, or disallowed method",
		})
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

	//
	// swis pkg registration
	//

	// alvax CRUD
	alvax.Cache = &core.Cache{}
	alvax.Routes(router.Group("/alvax"))

	// backups CRUD
	backups.Cache = &core.Cache{}
	backups.Routes(router.Group("/backups"))

	// business CRUD
	business.Cache = &core.Cache{}
	business.Routes(router.Group("/business"))

	// depots CRUD
	depots.Cache = &core.Cache{}
	depots.Routes(router.Group("/depots"))

	// dish CRUD
	dish.Dispatcher = dish.NewDispatcher()
	dish.CacheIncidents = &core.Cache{}
	dish.CacheSockets = &core.Cache{}
	dish.Routes(router.Group("/dish"))

	// finance accounts CRUD
	finance.CacheAccounts = &core.Cache{}
	finance.CacheItems = &core.Cache{}
	finance.Routes(router.Group("/finance"))

	// infra CRUD
	infra.CacheHosts = &core.Cache{}
	infra.CacheNetworks = &core.Cache{}
	infra.CacheDomains = &core.Cache{}
	infra.Routes(router.Group("/infra"))

	// links CRUD
	links.Cache = &core.Cache{}
	links.Routes(router.Group("/links"))

	// news CRUD
	news.Cache = &core.Cache{}
	news.Routes(router.Group("/news"))

	// projects CRUD
	projects.Cache = &core.Cache{}
	projects.Routes(router.Group("/projects"))

	// roles CRUD
	roles.Cache = &core.Cache{}
	roles.Routes(router.Group("/roles"))

	// users CRUD
	users.Cache = &core.Cache{}
	users.Routes(router.Group("/users"))

	// attach router to http.Server and start it, check for SERVER_PORT env variable
	if os.Getenv("SERVER_PORT") == "" {
		log.Fatal("SERVER_PORT environment variable not provided! refusing to start the server...")
	}

	// https://pkg.go.dev/net/http#Server
	server := &http.Server{
		Addr:        "0.0.0.0:" + os.Getenv("SERVER_PORT"),
		Handler:     router,
		ReadTimeout: 10 * time.Second,
		//WriteTimeout: 10 * time.Second,
		// = 1 * 2^23 = 1,048,576 * 8
		MaxHeaderBytes: 1 << 23,
		// use config.CORSMiddleware()
		DisableGeneralOptionsHandler: true,
	}
	server.ListenAndServe()
}
