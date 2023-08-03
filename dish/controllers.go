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

var Cache *config.Cache

// (GET /dish/sockets)
// @Summary Get all sockets list
// @Description get socket list, socket array
// @Tags dish
// @Produce  json
// @Success 200 {object} string "ok"
// @Router /dish/sockets [get]
// Get all sockets loaded.
func GetSocketList(c *gin.Context) {
	sockets, count := Cache.GetAll()

	/*rawJSON, err := json.Marshal(sockets)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "cannot marshal sockets into JSON byte stream",
		})
	}

	rawChecksum := sha256.Sum256([]byte(rawJSON))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "cannot calculate the checksum",
		})
	}*/

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"count":    count,
		"message":  "ok, dumping all sockets",
		//"checksum": fmt.Sprintf("%x", rawChecksum),
		"sockets":  sockets,
	})
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
func GetSocketListByHost(c *gin.Context) {
	var host string = c.Param("host")
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
		c.IndentedJSON(http.StatusOK, gin.H{
			"code":     http.StatusOK,
			"message":  "ok, dumping socket by host",
			"host":     host,
			"sockets":  exportedSockets,
		})
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "no sockets for given 'host'",
		"host":    host,
	})
	return

}

// @Summary Adding new socket to socket array
// @Description add new socket to socket array
// @Tags dish
// @Produce json
// @Param request body dish.Socket true "query params"
// @Success 200 {object} dish.Socket
// @Router /dish/sockets [post]
// Add new socket to the list.
func PostNewSocket(c *gin.Context) {
	var newSocket Socket

	if err := c.BindJSON(&newSocket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot bind input JSON stream",
		})
		return
	}

	if _, found := Cache.Get(newSocket.ID); found {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "dish socket already exists",
			"id":      newSocket.ID,
		})
		return
	}

	if saved := Cache.Set(newSocket.ID, newSocket); !saved {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "socket couldn't be saved to database",
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "socket added",
		"socket":  newSocket,
	})
	return
}

// (PUT /dish/sockets/{id})
// @Summary Update socket by its ID
// @Description update socket by its ID
// @Tags dish
// @Produce json
// @Param request body dish.Socket.ID true "query params"
// @Success 200 {object} dish.Socket
// @Router /dish/sockets/{id} [put]
// edit existing socket by ID
func UpdateSocketByID(c *gin.Context) {
	var id string = c.Param("id")
	var updatedSocket Socket

	if _, found := Cache.Get(id); !found {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "socket not found",
		})
		return
	}

	if err := c.BindJSON(&updatedSocket); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot bind input JSON stream",
		})
		return
	}

	if saved := Cache.Set(updatedSocket.ID, updatedSocket); !saved {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "socket couldn't be saved to database",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "socket updated by its ID",
		"socket":  updatedSocket,
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
// @Router /dish/sockets/{id}/mute [put]
// edit existing socket by ID
func MuteToggleSocketByID(c *gin.Context) {
	var id string = c.Param("id")
	var updatedSocket Socket

	rawSocket, ok := Cache.Get(id)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "socket not found",
			"id":      id,
		})
		return
	}

	updatedSocket, ok = rawSocket.(Socket)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "socket couldn't be saved to database",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "socket mute toggle pressed!",
		"socket":  updatedSocket,
	})
	return
}

// (DELETE /dish/sockets/{id})
// remove existing socket by ID
// @Summary Delete socket by its ID
// @Description delete socket by its ID
// @Tags dish
// @Produce json
// @Param  id  path  string  true  "dish ID"
// @Success 200 {object} dish.Socket
// @Router /dish/sockets/{id} [delete]
func DeleteSocketByID(c *gin.Context) {
	var id string = c.Param("id")

	if _, found := Cache.Get(id); !found {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "socket not found",
			"id":      id,
		})
		return
	}

	if deleted := Cache.Delete(id); !deleted {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "socket couldn't be deleted from database",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "socket deleted by ID",
		"id":      id,
	})
	return
}

// (POST /dish/sockets/restore)
// @Summary Upload dish dump backup -- restores all loaded sockets
// @Description update dish JSON dump
// @Tags dish
// @Accept json
// @Produce json
// @Router /dish/restore [post]
func PostDumpRestore(c *gin.Context) {
	var importSockets = &Sockets{}
	var socket Socket
	var counter int = 0

	if err := c.BindJSON(importSockets); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot bind input JSON stream",
		})
		return
	}

	for _, socket = range importSockets.Sockets {
		Cache.Set(socket.ID, socket)
		counter++
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"count":   counter,
		"message": "sockets imported, omitting output",
	})
	return
}
