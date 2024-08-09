package dish

import (
	"github.com/gin-gonic/gin"
)

func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetDishRoot)
	g.POST("/restore",
		PostDumpRestore)
	g.GET("/types",
		ListTypes)

	g.GET("/incidents",
		GetIncidentList)
	g.GET("/incidents/global",
		GetGlobalIncidentList)
	g.GET("/incidents/public",
		GetPublicIncidentList)
	g.GET("/incidents/:key",
		GetIncidentListBySocketID)
	g.POST("/incidents",
		PostNewIncident)
	g.PUT("/incidents/:key",
		UpdateIncidentByKey)
	g.PATCH("/incidents/:key",
		UpdateIncidentByKey)
	g.DELETE("/incidents/:key",
		DeleteIncidentByKey)

	g.GET("/sockets",
		GetSocketList)
	g.GET("/sockets/status",
		GetSSEvents)
	g.GET("/sockets/public",
		GetSocketListPublic)
	g.GET("/sockets/:host",
		GetSocketListByHost)
	g.POST("/sockets/:key",
		PostNewSocketByKey)
	g.PUT("/sockets/:key",
		UpdateSocketByKey)
	g.PATCH("/sockets/:key",
		UpdateSocketByKey)
	g.PUT("/sockets/:key/mute",
		MuteToggleSocketByKey)
	g.PUT("/sockets/:key/maintenance",
		MaintenanceToggleSocketByKey)
	g.PUT("/sockets/:key/public",
		PublicToggleSocketByKey)
	g.DELETE("/sockets/:key",
		DeleteSocketByKey)
	g.POST("/sockets/results",
		BatchPostHealthyStatus)
}
