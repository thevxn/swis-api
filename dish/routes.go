package dish

import (
	"github.com/gin-gonic/gin"
)

func Routes(g *gin.RouterGroup) {
	g.HEAD("/test",
		HeadTest)
	g.GET("/sockets",
		GetSocketList)
	g.GET("/sockets/:host",
		GetSocketListByHost)
	g.POST("/sockets",
		PostNewSocket)
	g.PUT("/sockets/:id",
		UpdateSocketByID)
	g.PATCH("/sockets/:id",
		UpdateSocketByID)
	g.PUT("/sockets/:id/mute",
		MuteToggleSocketByID)
	g.DELETE("/sockets/:id",
		DeleteSocketByID)
	g.POST("/sockets/restore",
		PostDumpRestore)
}
