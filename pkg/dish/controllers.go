package dish

import (
	"encoding/json"
	"io"
	//"log"
	"net/http"
	"time"

	//"go.savla.dev/dish/pkg/socket"
	"go.savla.dev/swis/v5/pkg/core"

	"github.com/gin-gonic/gin"
)

var (
	Cache   *core.Cache
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
	core.PrintAllRootItems(ctx, Cache, pkgName)
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
	core.AddNewItemByParam(ctx, Cache, pkgName, Socket{})
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
	core.UpdateItemByParam(ctx, Cache, pkgName, Socket{})
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
	core.DeleteItemByParam(ctx, Cache, pkgName)
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
	core.BatchRestoreItems(ctx, Cache, pkgName, Socket{})
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
	var counter int = 0

	rawSocketsMap, _ := Cache.GetAll()

	for _, rawSocket := range rawSocketsMap {
		socket, ok := rawSocket.(Socket)
		if !ok {
			continue
		}

		if contains(socket.DishTarget, host) && !socket.Muted {
			exportedSockets[socket.ID] = socket
			counter++
		}
	}

	if len(exportedSockets) > 0 {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"count":   counter,
			"items":   exportedSockets,
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

// @Summary Batch update socket's healthy state.
// @Description batch update socket's healthy state.
// @Tags dish
// @Produce json
// @Router /dish/sockets/results [post]
func BatchPostHealthyStatus(ctx *gin.Context) {
	var results = struct {
		Map map[string]bool `json:"dish_results"`
	}{}

	if err := ctx.BindJSON(&results); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
			"message": "cannot bind input JSON stream",
			"package": pkgName,
		})
		return
	}

	var sockets []string
	var count int = 0

	for key, result := range results.Map {
		var socket Socket
		var ok bool

		if rawSocket, found := Cache.Get(key); !found {
			continue
		} else {
			if socket, ok = rawSocket.(Socket); !ok {
				ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"key":     key,
					"message": "cannot assert data type, database internal error",
					"package": pkgName,
				})
				return
			}
		}

		socket.Healthy = result
		socket.TestTimestamp = time.Now().UnixNano()

		sockets = append(sockets, socket.ID)
		count++

		if saved := Cache.Set(key, socket); !saved {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"key":     key,
				"message": "cannot update socket's healthy state by key",
			})
			return
		}
	}

	// generate and send a SSE message
	msg := Message{
		Content:    "sockets updated",
		SocketList: sockets,
		Timestamp:  time.Now().UnixNano(),
	}

	Dispatcher.NewEvent(msg)

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, healthy booleans updated per socket",
		"count":   count,
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

// @Summary Subscribe to dish SSE service.
// @Description subscribe to dish SSE service
// @Tags dish
// @Produce json
// @Success 200 {object} dish.Message
// @Router /dish/sockets/status [get]
func SubscribeToSSEStream(ctx *gin.Context) {
	// initialize client channel
	clientChan := make(ClientChan)

	// send new connection to event server
	Dispatcher.NewClients <- clientChan

	defer func() {
		// send closed connection to event server
		Dispatcher.ClosedClients <- clientChan
	}()

	// set the stream headers
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")

	ctx.Stream(func(w io.Writer) bool {
		// Stream message to client from message channel
		if msg, ok := <-clientChan; ok {
			m, _ := json.Marshal(msg)
			ctx.SSEvent("message", string(m))
			//log.Println("wrote:", m)
			return true
		}
		return false
	})
}
