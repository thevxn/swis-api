// Package swis-api is RESTful API core backend aka 'sakalWeb Information System v5'.
// Basically it is a system of high modularity, where each module (package in golang terminology)
// has its routes, models, and controllers (handler functions) defined in its own folder.
package main

import (
	// golang libs

	"net/http"
	"time"

	// swapi modules -- very local dependencies
	"swis-api/alvax"
	"swis-api/auth"
	//"swis-api/auth"
	"swis-api/business"
	"swis-api/depot"
	"swis-api/dish"
	"swis-api/docs"
	"swis-api/finance"
	"swis-api/groups"
	"swis-api/infra"
	"swis-api/news"
	"swis-api/projects"
	"swis-api/swife"
	"swis-api/users"

	// remote dependencies

	gin "github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title swis-api v5
// @version 5.2
// @description sakalWeb Information System v5 RESTful API documentation
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://savla.dev/swapi
// @contact.email info@savla.dev

// @license.name MIT
// @license.url https://github.com/savla-dev/swis-api/blob/master/LICENSE

// @host swapi.savla.su
// @BasePath /
func main() {
	// blank gin without any middleware
	router := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// by default gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// use custom swapi Auth middleware
	//router.Use(auth.SwapiAuth())

	// JWT middleware
	/*
		authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
			Realm:           "swapi-dev-zone",
			Key:             []byte("sekret tve mamy"),
			Timeout:         time.Hour,
			MaxRefresh:      time.Hour,
			IdentityKey:     auth.IdentityKey,
			PayloadFunc:     auth.PayloadFunc,
			IdentityHandler: auth.IdentityHandler,
			Authenticator:   auth.Authenticator,
			Authorizator:    auth.Authorizator,
			Unauthorized:    auth.Unauthorized,

			// TokenLookup is a string in the form of "<source>:<name>" that is used
			// to extract token from the request.
			// Optional. Default value "header:Authorization".
			// Possible values:
			// - "header:<name>"
			TokenLookup: "header:<name>",
			// - "query:<name>"
			// - "cookie:<name>"
			// - "param:<name>"
			//TokenLookup: "header: Authorization, query: token, cookie: jwt",
			// TokenLookup: "query:token",
			// TokenLookup: "cookie:token",

			// TokenHeadName is a string in the header. Default value is "Bearer"
			TokenHeadName: "Bearer",

			// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
			TimeFunc: time.Now,
		})

		if err != nil {
			log.Fatal("JWT Error:" + err.Error())
		}

		// When you use jwt.New(), the function is already automatically called for checking,
		// which means you don't need to call it again.
		errInit := authMiddleware.MiddlewareInit()

		if errInit != nil {
			log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
		}*/

	//router.GET("/refresh_token", authMiddleware.RefreshHandler)
	//router.POST("/login", authMiddleware.LoginHandler)

	// root path --- auth required
	router.GET("/", func(c *gin.Context) {
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

	// swagger documentation
	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://0.0.0.0:8049/docs/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
