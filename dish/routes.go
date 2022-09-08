package dish

import (
	"github.com/gin-gonic/gin"
)

func Routes(g *gin.RouterGroup) {
	// @Summary
	// @Description
	// @Tags dish
	// @Success 200
	// @Router /dish/test [head]
	g.HEAD("/test", HeadTest)

	// @Summary Get all sockets list
	// @Description get socket list, socket array
	// @Tags dish
	// @Produce  json
	// @Success 200 {object} string "ok"
	// @Router /dish/sockets [get]
	g.GET("/sockets", GetSocketList)

	// @Summary Get socket list by host
	// @Description get socket list by Host
	// @Tags dish
	// @Produce  json
	// @Param   host     path    string     true        "dish instance name"
	// @Success 200 {string} string	"ok"
	// @Router /dish/sockets/{host} [get]
	g.GET("/sockets/:host", GetSocketListByHost)

	// @Summary Adding new socket to socket array
	// @Description add new socket to socket array
	// @Tags dish
	// @Produce json
	// @Param request body dish.Socket true "query params"
	// @Success 200 {object} dish.Socket
	g.POST("/sockets", PostNewSocket)

	// @Summary Update socket by its ID
	// @Description update socket by its ID
	// @Tags dish
	// @Produce json
	// @Param request body dish.Socket.ID true "query params"
	// @Success 200 {object} dish.Socket
	// @Router /dish/sockets/{id} [put]
	g.PUT("/sockets/:id", UpdateSocketByID)

	// (PATCH /sockets/{id})
	// edit existing socket by ID
	g.PATCH("/sockets/:id", UpdateSocketByID)

	// @Summary Mute/unmute socket by its ID
	// @Description mute/unmute socket by its ID
	// @Tags dish
	// @Produce json
	// @Param  id  path  string  true  "dish ID"
	// @Success 200 {object} dish.Socket
	// @Router /dish/sockets/{id}/mute [put]
	g.PUT("/sockets/:id/mute", MuteToggleSocketByID)

	// @Summary Delete socket by its ID
	// @Description delete socket by its ID
	// @Tags dish
	// @Produce json
	// @Success 200 {string} string "ok"
	// @Router /dish/sockets/{id} [delete]
	g.DELETE("/sockets/:id", DeleteSocketByID)

	// @Summary Upload dish dump backup -- restores all loaded sockets
	// @Description update dish JSON dump
	// @Tags dish
	// @Accept json
	// @Produce json
	// @Router /dish/restore [post]
	g.POST("/sockets/restore", PostDumpRestore)
}
