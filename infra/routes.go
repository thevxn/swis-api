package infra

import (
	"github.com/gin-gonic/gin"
)

// infra CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetInfrastructure)
	g.GET("/hosts",
		GetHosts)
	g.GET("/hosts/:hostname",
		GetHostByHostname)
	//g.GET("/map", GetInfraMap)
	g.GET("/networks",
		GetNetworks)
	g.POST("/restore",
		PostDumpRestore)
}
