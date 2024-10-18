package dish

import (
	"github.com/gin-gonic/gin"
)

func Routes(g *gin.RouterGroup) {
	// common
	g.GET("",
		GetDishRoot)
	g.POST("/restore",
		PostDumpRestore)

	// incidents
	g.GET("/incidents",
		GetIncidentList)
	g.POST("/incidents",
		PostNewIncident)
	g.GET("/incidents/types",
		ListTypesIncidents)
	g.GET("/incidents/global",
		GetGlobalIncidentList)
	g.GET("/incidents/public",
		GetPublicIncidentList)
	g.GET("/incidents/:key",
		GetIncidentListBySocketID)
	g.PUT("/incidents/:key",
		UpdateIncidentByKey)
	g.PATCH("/incidents/:key",
		UpdateIncidentByKey)
	g.DELETE("/incidents/:key",
		DeleteIncidentByKey)

	// sockets
	g.GET("/sockets",
		GetSocketList)
	g.POST("/sockets",
		PostNewSocket)
	g.GET("/sockets/types",
		ListTypesSockets)
	g.GET("/sockets/status",
		GetSSEvents)
	g.GET("/sockets/public",
		GetSocketListPublic)
	g.GET("/sockets/:host",
		GetSocketListByHost)
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

	// streamer stats
	g.GET("/streamer/stats",
		GetStreamerStats)
}
