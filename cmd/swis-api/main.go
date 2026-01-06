// @title swis-api (swapi) v5
// @version 5.19.0
// @description sakalWeb Information System v5 RESTful API documentation
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://vxn.dev/swapi
// @contact.email krusty@vxn.dev

// @license.name MIT
// @license.url https://github.com/vxn-dev/swis-api/blob/master/LICENSE

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

func main() {
	server := newServer()
	server.Run()
}
