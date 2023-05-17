package dish

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/savla-dev/savla-dish/socket"
)

var s sync.Map

// (HEAD /dish/test)
// @Summary
// @Description
// @Tags dish
// @Success 200
// @Router /dish/test [head]
// HeadTest is the HEAD HTTP method for savla-dish service, that acts like a testing endpoint.
func HeadTest(c *gin.Context) {
	c.String(http.StatusOK, "ok")
	return
}

// (GET /dish/sockets)
// @Summary Get all sockets list
// @Description get socket list, socket array
// @Tags dish
// @Produce  json
// @Success 200 {object} string "ok"
// @Router /dish/sockets [get]
// Get all sockets loaded.
func GetSocketList(c *gin.Context) {
	var sockets = make(map[string]Socket)

	s.Range(func(rawKey, rawVal interface{}) bool {
		// very insecure assert?
		k, ok := rawKey.(string)
		v, ok := rawVal.(Socket)

		if !ok {
			return false
		}

		sockets[k] = v
		return true
	})

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, dumping all sockets",
		"sockets": sockets,
	})
	return
}

// (GET /dish/sockets/{host})
// @Summary Get socket list by host
// @Description get socket list by Host
// @Tags dish
// @Produce  json
// @Param   host     path    string     true        "dish instance name"
// @Success 200 {string} string	"ok"
// @Router /dish/sockets/{host} [get]
// Get sockets by hostname/dish-name.
func GetSocketListByHost(c *gin.Context) {
	var host string = c.Param("host")
	var sockets = make(map[string]Socket)

	s.Range(func(rawKey, rawVal interface{}) bool {
		// very insecure assert?
		k, ok := rawKey.(string)
		v, ok := rawVal.(Socket)

		if !ok {
			return false
		}

		if contains(v.DishTarget, host) && !v.Muted {
			sockets[k] = v
		}
		return true
	})

	if len(sockets) > 0 {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "ok, dumping socket by host",
			"code":    http.StatusOK,
			"sockets": sockets,
		})
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{
		"message": "no sockets for given 'hostname'",
		"code":    http.StatusNotFound,
		"host":    host,
	})
	return

}

// (POST /sockets)
// @Summary Adding new socket to socket array
// @Description add new socket to socket array
// @Tags dish
// @Produce json
// @Param request body dish.Socket true "query params"
// @Success 200 {object} dish.Socket
// Add new socket to the list.
func PostNewSocket(c *gin.Context) {
	var newSocket = &Socket{}

	// bind JSON to newSocket
	if err := c.BindJSON(newSocket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	if _, found := s.Load(newSocket.ID); found {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"message": "dish socket already exists",
			"code":    http.StatusConflict,
			"id":      newSocket.ID,
		})
		return
	}

	s.Store(newSocket.ID, newSocket)

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
	var updatedSocket = &Socket{}

	if _, ok := s.Load(id); !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "socket not found",
		})
		return
	}

	if err := c.BindJSON(updatedSocket); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	s.Store(updatedSocket.ID, updatedSocket)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "socket updated",
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
	var updatedSocket Socket
	var id string = c.Param("id")

	rawSocket, ok := s.Load(id)
	updatedSocket, ok = rawSocket.(Socket)

	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "socket not found",
			"id":      id,
		})
		return
	}

	// inverse the Muted field value
	updatedSocket.Muted = !updatedSocket.Muted
	if updatedSocket.Muted {
		updatedSocket.MutedFrom = time.Now().Unix()
	}

	s.Store(updatedSocket.ID, updatedSocket)

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

	if _, ok := s.Load(id); !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "socket not found",
			"id":      id,
		})
		return
	}

	s.Delete(id)

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

	if err := c.BindJSON(importSockets); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	for _, socket = range importSockets.Sockets {
		s.Store(socket.ID, socket)
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "sockets imported, omitting output",
	})
	return
}
