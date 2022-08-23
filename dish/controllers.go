package dish

import (
	"net/http"

	"github.com/gin-gonic/gin"
	//"github.com/savla-dev/savla-dish/socket"
)

// @Summary
// @Description
// @Tags dish
// @Success 200
// @Router /dish/test [head]
// HeadTest is the HEAD HTTP method for savla-dish service, that acts like a testing endpoint
func HeadTest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"test": true,
	})
}

// @Summary Get all sockets list
// @Description get socket list, socket array
// @Tags dish
// @Produce  json
// @Success 200 {object} string "ok"
// @Router /dish/sockets [get]
// GetSocketList GET method
func GetSocketList(c *gin.Context) {
	//var sockets = socket.Sockets{
	var sockets = Sockets{
		Sockets: socketArray,
	}

	c.IndentedJSON(http.StatusOK, sockets)
}

// @Summary Get socket list by host
// @Description get socket list by Host
// @Tags dish
// @Produce  json
// @Param   host     path    string     true        "dish instance name"
// @Success 200 {string} string	"ok"
// @Router /dish/sockets/{host} [get]
// GetSocketListByHost GET
func GetSocketListByHost(c *gin.Context) {
	host := c.Param("host")

	//var sockets = socket.Sockets{
	var sockets = Sockets{}
	/*var sockets = Sockets{
		Sockets: []Socket{},
	}*/

	// loop over socketArray, find
	for _, s := range socketArray {
		if contains(s.DishTarget, host) {
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

// @Summary Adding new socket to socket array
// @Description add new socket to socket array
// @Tags dish
// @Produce json
// @Param request body dish.Socket true "query params"
// @Success 200 {object} dish.Socket
// @Router /dish/sockets [post]
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

// @Summary Update socket by its ID
// @Description update socket by its ID
// @Tags dish
// @Produce json
// @Param request body dish.Socket.ID true "query params"
// @Success 200 {object} dish.Socket
// @Router /dish/sockets/{id} [put]
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

// @Summary Delete socket by its ID
// @Description delete socket by its ID
// @Tags dish
// @Produce json
// @Success 200 {string} string "ok"
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

// @Summary Upload dish dump backup -- restores all loaded sockets
// @Description update dish JSON dump
// @Tags dish
// @Accept json
// @Produce json
// @Router /dish/restore [post]
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
