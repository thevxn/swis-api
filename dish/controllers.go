package dish

import (
	"net/http"

	"github.com/gin-gonic/gin"
	//"github.com/savla-dev/savla-dish/socket"
)

// PostSocketTestResult
/*
func PostSocketTestResult(c *gin.Context) {
	var importResults Sockets

	if err := c.BindJSON(&importSockets); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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
*/

// HeadTest is the HEAD HTTP method for savla-dish service, that acts like a testing endpoint
func HeadTest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"test": true,
	})
}

// GetSocketList GET method
func GetSocketList(c *gin.Context) {
	//var sockets = socket.Sockets{
	var sockets = Sockets{
		Sockets: socketArray,
	}

	c.IndentedJSON(http.StatusOK, sockets)
}

// contains checks if a string is present in a slice
// https://freshman.tech/snippets/go/check-if-slice-contains-element/
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

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
		if contains(s.DishList, host) {
			// clear the dish source list for the client (dish)
			s.DishList = []string{host}
			//s.DishList = nil
			sockets.Sockets = append(sockets.Sockets, s)
		}
	}

	if len(sockets.Sockets) > 0 {
		c.IndentedJSON(http.StatusOK, sockets)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no sockets for given 'hostname'"})

}

func PostDumpRestore(c *gin.Context) {
	var importSockets Sockets

	if err := c.BindJSON(&importSockets); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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
