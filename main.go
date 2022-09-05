// @title swis-api v5
// @version 5.2
// @description sakalWeb Information System v5 RESTful API documentation
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://savla.dev/swapi
// @contact.email info@savla.dev

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
	"fmt"
	"net/http"
	"time"

	// swapi modules -- very local dependencies
	"swis-api/alvax"
	"swis-api/auth"

	//"swis-api/auth"
	"swis-api/business"
	"swis-api/depot"
	"swis-api/dish"
	"swis-api/finance"
	"swis-api/infra"
	"swis-api/news"
	"swis-api/projects"
	"swis-api/roles"
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

	// use custom swapi Auth middleware --- token auth
	router.Use(auth.AuthMiddleware())

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// by default gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// custom logging format
		return fmt.Sprintf("%s (%s) [%s] | %s %s | %s | %d | %s | \"%s\" %s\"\n",
			param.ClientIP,
			auth.Params.User.Name,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// root path
	router.GET("/", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{
			"title":   "sakalWebIS v5 RESTful API -- root route",
			"message": "welcome to swis, " + auth.Params.User.Name + "!",
			"code":    http.StatusOK,
			//"bearer":    auth.Params.BearerToken,
			"timestamp": time.Now().Unix(),
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

	// infra CRUD
	infraRouter := router.Group("/infra")
	infra.Routes(infraRouter)

	// news CRUD
	newsRouter := router.Group("/news")
	news.Routes(newsRouter)

	// projects CRUD
	projectsRouter := router.Group("/projects")
	projects.Routes(projectsRouter)

	// roles CRUD
	rolesRouter := router.Group("/roles")
	roles.Routes(rolesRouter)

	// swife CRUD
	swifeRouter := router.Group("/swife")
	swife.Routes(swifeRouter)

	// users CRUD
	usersRouter := router.Group("/users")
	//usersRouter.Use(authMiddleware.MiddlewareFunc())
	users.Routes(usersRouter)

	// webui CRUD
	// only a small proposal for simple swapi fronted to browse swapi data quickly
	// import "swis-api/webui"
	//webuiRouter := router.Group("/webui")
	//webui.Routes(webuiRouter)

	// attach router to http.Server and start it
	// https://pkg.go.dev/net/http#Server
	server := &http.Server{
		Addr:         "0.0.0.0:8049",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		//MaxHandlerBytes: 1 << 20,
	}
	server.ListenAndServe()
}
