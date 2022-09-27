package dish

import (
	"net/http"
	"swis-api/infra"

	"github.com/gin-gonic/gin"
	//"github.com/savla-dev/savla-dish/socket"
)

var socketArray = []Socket{}

// (HEAD /dish/test)
// @Summary
// @Description
// @Tags dish
// @Success 200
// @Router /dish/test [head]
// HeadTest is the HEAD HTTP method for savla-dish service, that acts like a testing endpoint.
func HeadTest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"test": true,
	})
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
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, dumping all sockets",
		"sockets": socketArray,
	})
}

// (GET /dish/targets)
// GetTargetList GET method
func GetTargetList(c *gin.Context) {
	var targets = infra.Hosts{}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, dumping all targets",
		"targets": targets,
	})
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
	host := c.Param("host")

	//var sockets = socket.Sockets{
	var sockets = Sockets{}
	/*var sockets = Sockets{
		Sockets: []Socket{},
	}*/

	// loop over socketArray, find
	for _, s := range socketArray {
		// export only unmuted sockets
		if contains(s.DishTarget, host) && !s.Muted {
			// clear the dish source list for the client (dish)
			s.DishTarget = []string{host}
			sockets.Sockets = append(sockets.Sockets, s)
		}
	}

	if len(sockets.Sockets) > 0 {
		c.IndentedJSON(http.StatusOK, sockets)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no sockets for given 'hostname'"})

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
	var newSocket Socket

	// bind JSON to newSocket
	if err := c.BindJSON(&newSocket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	/*
		if _, s := findSocketByID(c); s.ID == newSocket.ID {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "this socket ID already exists! not allowed to POST",
			})
		}
	*/

	// add new socket
	socketArray = append(socketArray, newSocket)

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "socket added",
		"socket":  newSocket,
	})
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
	var updatedSocket Socket

	i, _ := findSocketByID(c)

	if err := c.BindJSON(&updatedSocket); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	socketArray[*i] = updatedSocket
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

	i, _ := findSocketByID(c)
	updatedSocket = socketArray[*i]

	// inverse the Muted field value
	updatedSocket.Muted = !updatedSocket.Muted

	socketArray[*i] = updatedSocket
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
	i, s := findSocketByID(c)

	// delete an element from the array
	// https://www.educative.io/answers/how-to-delete-an-element-from-an-array-in-golang
	newLength := 0
	for index := range socketArray {
		if *i != index {
			socketArray[newLength] = socketArray[index]
			newLength++
		}
	}

	// reslice the array to remove extra index
	socketArray = socketArray[:newLength]

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "socket deleted by ID",
		"socket":  *s,
	})
}

// (POST /dish/sockets/restore)
// @Summary Upload dish dump backup -- restores all loaded sockets
// @Description update dish JSON dump
// @Tags dish
// @Accept json
// @Produce json
// @Router /dish/restore [post]
// restore all sockets from JSON dump (JSON-bind)
func PostDumpRestore(c *gin.Context) {
	var importSockets Sockets

	if err := c.BindJSON(&importSockets); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	//depots = append(depots, importDepot)
	socketArray = importSockets.Sockets

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "sockets imported, omitting output",
	})
}
