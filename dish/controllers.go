package dish

import (
	//"crypto/sha256"
	//"encoding/json"
	//"fmt"
	"net/http"
	"time"

	//"go.savla.dev/dish/pkg/socket"
	"go.savla.dev/swis/v5/config"

	"github.com/gin-gonic/gin"
)

var (
	Cache   *config.Cache
	pkgName string = "dish"
)

// Get all sockets loaded.
// @Summary Get all sockets list
// @Description get socket list, socket array
// @Tags dish
// @Produce  json
// @Success 200 {object} string "ok"
// @Router /dish/sockets [get]
func GetSocketList(ctx *gin.Context) {
	config.PrintAllRootItems(ctx, Cache, pkgName)
	return
}

// Add new socket to the list.
// @Summary Adding new socket to socket array
// @Description add new socket to socket array
// @Tags dish
// @Produce json
// @Param request body dish.Socket true "query params"
// @Success 200 {object} dish.Socket
// @Router /dish/sockets/{key} [post]
func PostNewSocketByKey(ctx *gin.Context) {
	config.AddNewItemByParam(ctx, Cache, pkgName, Socket{})
	return
}

// edit existing socket by ID
// @Summary Update socket by its ID
// @Description update socket by its ID
// @Tags dish
// @Produce json
// @Param request body dish.Socket.ID true "query params"
// @Success 200 {object} dish.Socket
// @Router /dish/sockets/{key} [put]
func UpdateSocketByKey(ctx *gin.Context) {
	config.UpdateItemByParam(ctx, Cache, pkgName, Socket{})
	return
}

// remove existing socket by ID
// @Summary Delete socket by its ID
// @Description delete socket by its ID
// @Tags dish
// @Produce json
// @Param  id  path  string  true  "dish ID"
// @Success 200 {object} dish.Socket
// @Router /dish/sockets/{key} [delete]
func DeleteSocketByKey(ctx *gin.Context) {
	config.DeleteItemByParam(ctx, Cache, pkgName)
	return
}

// @Summary Upload dish dump backup -- restores all loaded sockets
// @Description update dish JSON dump
// @Tags dish
// @Accept json
// @Produce json
// @Success 201
// @Router /dish/restore [post]
func PostDumpRestore(ctx *gin.Context) {
	config.BatchRestoreItems(ctx, Cache, pkgName, Socket{})
	return
}

// @Summary Get socket list by host
// @Description get socket list by Host
// @Tags dish
// @Produce json
// @Param host path string true "dish instance name"
// @Success 200 {string} string	"ok"
// @Router /dish/sockets/{host} [get]
// Get sockets by hostname/dish-name.
func GetSocketListByHost(ctx *gin.Context) {
	var host string = ctx.Param("host")
	var exportedSockets = make(map[string]Socket)

	rawSocketsMap, _ := Cache.GetAll()

	for _, rawSocket := range rawSocketsMap {
		socket, ok := rawSocket.(Socket)
		if !ok {
			continue
		}

		if contains(socket.DishTarget, host) && !socket.Muted {
			exportedSockets[socket.ID] = socket
		}
	}

	if len(exportedSockets) > 0 {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"tems":    exportedSockets,
			"message": "ok, dumping socket by host",
			"host":    host,
		})
		return
	}

	ctx.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "no sockets for given 'host'",
		"host":    host,
	})
	return

}

// (PUT /dish/sockets/{id}/mute)
// @Summary Mute/unmute socket by its ID
// @Description mute/unmute socket by its ID
// @Tags dish
// @Produce json
// @Param  id  path  string  true  "dish ID"
// @Success 200 {object} dish.Socket
// @Router /dish/sockets/{key}/mute [put]
// edit existing socket by ID
func MuteToggleSocketByKey(ctx *gin.Context) {
	var id string = ctx.Param("key")
	var updatedSocket Socket

	rawSocket, ok := Cache.Get(id)
	if !ok {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "socket not found",
			"id":      id,
		})
		return
	}

	updatedSocket, ok = rawSocket.(Socket)
	if !ok {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "cannot assert data type, database internal error",
		})
		return
	}

	// inverse the Muted field value
	updatedSocket.Muted = !updatedSocket.Muted

	if updatedSocket.Muted {
		updatedSocket.MutedFrom = time.Now().Unix()
	}

	if saved := Cache.Set(updatedSocket.ID, updatedSocket); !saved {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "socket couldn't be saved to database",
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "socket mute toggle pressed!",
		"socket":  updatedSocket,
	})
	return
}
