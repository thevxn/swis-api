package dish

import (
	"net/http"
	"swis-api/infra"

	"github.com/gin-gonic/gin"
	//"github.com/savla-dev/savla-dish/socket"
)

// (HEAD /test)
// HeadTest is the HEAD HTTP method for savla-dish service, that acts like a testing endpoint.
func HeadTest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"test": true,
	})
}

// (GET /sockets)
// Get all sockets loaded.
func GetSocketList(c *gin.Context) {
	//var sockets = socket.Sockets{
	var sockets = Sockets{
		Sockets: socketArray,
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, dumping all sockets",
		"sockets": sockets,
	})
}

// GetTargetList GET method
func GetTargetList(c *gin.Context) {
	var targets = infra.Hosts{}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, dumping all targets",
		"targets": targets,
	})
}

// (GET /sockets/{host})
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

// (PUT /sockets/{id}/mute)
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

// (PUT /sockets/{id})
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

// (DELETE /sockets/{id})
// remove existing socket by ID
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

// (POST /sockets/restore)
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
