package alvax

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)


type AlvaxCommands struct {
	User 		string 		`json:"user"`
	CommandList	[]Command	`json:"command_list"`
}

type Command struct {
	// name as in the '/name' Telegram command syntax
	Name		string		`json:"name"`
	AliasNames	[]string	`json:"alias_names"`
	ArgumentList	[]string	`json:"argument_list"`
	ParentClass	string		`json:"parent_class"`
	RequiredArg	bool		`json:"required_argument" default:false`
}


var commandList = []Command{
	{Name: "bomb", ParentClass: "Bomb",  ArgumentList: []string{"red", "green", "blue"}, RequiredArg: false},
	{Name: "dish", ParentClass: "Dish",  ArgumentList: []string{"enable", "disable", "mute", "search"}, RequiredArg: true},
	{Name: "kanban", ParentClass: "Kanban", ArgumentList: []string{"getAllProjects"}, RequiredArg: true},
	{Name: "memes", ParentClass: "Memes",  ArgumentList: []string{"megamind", "chad"}, RequiredArg: true},
	{Name: "rating", ParentClass: "Rating", ArgumentList: []string{"good", "bad"}, AliasNames: []string{"badbot", "goodbot"}},
}


// GetCommandList returns JSON serialized list of commands for the alvax backend.
func GetCommandList(c *gin.Context) {
	// serialize struct to JSON
	c.IndentedJSON(http.StatusOK, gin.H{
		"command_list": commandList,
	})
}


// PosCommandsDumpRestore
func PostCommandsDumpRestore(c *gin.Context) {
	var importCommands AlvaxCommands

	// bind received JSON to importCommands
	if err := c.BindJSON(&importCommands); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	commandList = importCommands.CommandList

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated,
		"message": "alvax command list imported successfully",
	})
}

