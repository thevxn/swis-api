package alvax

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get all alvax's commands
// @Description get command array for alvax
// @Tags alvax
// @Produce  json
// @Success 200 {object} alvax.AlvaxCommands
// @Router /alvax/commands [get]
// GetSocketList GET method
// GetCommandList returns JSON serialized list of commands for the alvax backend.
func GetCommandList(c *gin.Context) {
	// serialize struct to JSON
	c.IndentedJSON(http.StatusOK, gin.H{
		"command_list": commandList,
	})
}

// @Summary Upload alvax dump backup -- restores all loaded commands
// @Description update alvax JSON dump
// @Tags alvax
// @Accept json
// @Produce json
// @Router /alvax/commands/restore [post]
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
