package dish

import (
	"github.com/gin-gonic/gin"
)

func Routes(g *gin.RouterGroup) {
	g.GET("/sockets",
		GetSocketList)
	g.GET("/sockets/:host",
		GetSocketListByHost)
	g.POST("/sockets",
		PostNewSocketByKey)
	g.PUT("/sockets/:key",
		UpdateSocketByKey)
	g.PATCH("/sockets/:key",
		UpdateSocketByKey)
	g.PUT("/sockets/:key/mute",
		MuteToggleSocketByKey)
	g.DELETE("/sockets/:key",
		DeleteSocketByKey)
	g.POST("/sockets/restore",
		PostDumpRestore)
}
