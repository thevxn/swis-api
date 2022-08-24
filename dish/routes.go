package dish

import (
	"github.com/gin-gonic/gin"
)

func Routes(g *gin.RouterGroup) {
	// (HEAD /test)
	// testing route for dish
	g.HEAD("/test", HeadTest)

	// (GET /sockets)
	// get all sockets loaded
	g.GET("/sockets", GetSocketList)

	// (GET /sockets/{host})
	// get sockets by hostname/dish-name
	g.GET("/sockets/:host", GetSocketListByHost)

	// (POST /sockets)
	// add new socket to the list
	g.POST("/sockets", PostNewSocket)

	// (PUT /sockets/{id})
	// edit existing socket by ID
	g.PUT("/sockets/:id", UpdateSocketByID)

	// (PUT /sockets/{id}/mute)
	// edit existing socket by ID
	g.PUT("/sockets/:id/mute", MuteToggleSocketByID)

	// (PATCH /sockets/{id})
	// edit existing socket by ID
	g.PATCH("/sockets/:id", UpdateSocketByID)

	// (DELETE /sockets/{id})
	// remove existing socket by ID
	g.DELETE("/sockets/:id", DeleteSocketByID)

	// (POST /sockets/restore)
	// restore all sockets from JSON dump (JSON-bind)
	g.POST("/sockets/restore", PostDumpRestore)
}
