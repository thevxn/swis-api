package depot

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)


type Depot struct {
	Depot []Item `json:"users"`
}

type Item struct {
	ID       	string `json:"id"`
	Nickname 	string `json:"nickname"`
	Role     	string `json:"role"`
	TokenBase64	string `json:"tokenbase64"`
}

// users demo data for user struct
var users = []Item{
	{ID: "1", Nickname: "sysadmin", Role: "admin"},
	{ID: "2", Nickname: "dev", Role: "developer"},
	{ID: "3", Nickname: "op", Role: "operator"},
}

// GetDepot returns JSON serialized list of users and their properties.
func GetDepot(c *gin.Context) {
	// serialize struct to JSON
	// net/http response code
	c.IndentedJSON(http.StatusOK, users)
}

// GetItemByID returns user's properties, given sent ID exists in database.
func GetItemByID(c *gin.Context) {
	id := c.Param("id")

	// loop over users
	for _, a := range users {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

// PostItem enables one to add new user to users model.
func PostItem(c *gin.Context) {
	var newItem Item

	// bind received JSON to newItem
	if err := c.BindJSON(&newItem); err != nil {
		return
	}

	// add new user
	users = append(users, newItem)
	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, newItem)
}

