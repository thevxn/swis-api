package backups

import (
	"github.com/gin-gonic/gin"
)

func Routes(g *gin.RouterGroup) {
	// @Summary Get all backups status
	// @Description get backups actual status
	// @Tags backup
	// @Produce  json
	// @Success 200 {object} string "ok"
	// @Router /backups [get]
	g.GET("/backups", GetSocketList)

	// @Summary Get backup status by project/service
	// @Description get backup status by project/service
	// @Tags backups
	// @Produce  json
	// @Param   host     path    string     true        "dish instance name"
	// @Success 200 {string} string	"ok"
	// @Router /dish/sockets/{host} [get]
	g.GET("/backups/:project", GetSocketListByHost)

	// @Summary Adding new socket to socket array
	// @Description add new socket to socket array
	// @Tags backups
	// @Produce json
	// @Param request body dish.Socket true "query params"
	// @Success 200 {object} dish.Socket
	g.POST("/backups", PostNewSocket)

	// @Summary Update socket by its ID
	// @Description update socket by its ID
	// @Tags dish
	// @Produce json
	// @Param request body dish.Socket.ID true "query params"
	// @Success 200 {object} dish.Socket
	// @Router /dish/sockets/{id} [put]
	g.PUT("/sockets/:id", UpdateSocketByID)

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
