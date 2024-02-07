package infra

import (
	"github.com/gin-gonic/gin"
)

// infra CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetInfrastructure)
	g.POST("/restore",
		PostDumpRestore)

	// domains CRUD
	g.GET("/domains",
		GetDomains)
	g.GET("/domains/:key",
		GetDomainByKey)
	g.POST("/domains/:key",
		PostNewDomainByKey)
	g.PUT("/domains/:key",
		UpdateDomainByKey)
	g.DELETE("/domains/:key",
		DeleteDomainByKey)

	// hosts CRUD
	g.GET("/hosts",
		GetHosts)
	g.GET("/hosts/:key",
		GetHostByKey)
	g.POST("/hosts/:key/config",
		PostHostConfigByKey)
	g.POST("/hosts/:key/facts",
		PostHostFactsByKey)
	g.POST("/hosts/:key/vmic",
		PostHostVMICByKey)
	g.DELETE("/hosts/:key/vmic/:vm",
		DeleteHostVMICByKeyAndVM)
	g.POST("/hosts/:key",
		PostNewHostByKey)
	g.PUT("/hosts/:key",
		UpdateHostByKey)
	g.DELETE("/hosts/:key",
		DeleteHostByKey)

	// networks CRUD
	g.GET("/networks",
		GetNetworks)
	g.GET("/networks/:key",
		GetNetworkByKey)
	g.POST("/networks/:key",
		PostNewNetworkByKey)
	g.PUT("/networks/:key",
		UpdateNetworkByKey)
	g.DELETE("/networks/:key",
		DeleteNetworkByKey)
}
