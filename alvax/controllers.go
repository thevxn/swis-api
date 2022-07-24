package alvax

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCommandList returns JSON serialized list of commands for the alvax backend.
func GetCommandList(c *gin.Context) {
	// serialize struct to JSON
	c.IndentedJSON(http.StatusOK, gin.H{
		"command_list": commandList,
	})
}

// PosCommandsDumpRestore
func PostDumpRestore(c *gin.Context) {
	var importCommands AlvaxCommands

	// bind received JSON to importCommands
	if err := c.BindJSON(&importCommands); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	commandList = importCommands.CommandList

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "alvax command list imported successfully",
	})
}
