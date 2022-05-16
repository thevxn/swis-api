package users

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)


type Users struct {
	Users []User `json:"users"`
}

type User struct {
	ID       	string `json:"id"`
	Nickname 	string `json:"nickname"`
	Role     	string `json:"role"`
	TokenBase64	string `json:"tokenbase64"`
}

// users demo data for user struct
var users = []User{
	{ID: "1", Nickname: "sysadmin", Role: "admin"},
	{ID: "2", Nickname: "dev", Role: "developer"},
	{ID: "3", Nickname: "op", Role: "operator"},
}

// GetUsers returns JSON serialized list of users and their properties.
func GetUsers(c *gin.Context) {
	// serialize struct to JSON
	// net/http response code
	c.IndentedJSON(http.StatusOK, users)
}

// GetUserByID returns user's properties, given sent ID exists in database.
func GetUserByID(c *gin.Context) {
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

// PostUser enables one to add new user to users model.
func PostUser(c *gin.Context) {
	var newUser User

	// bind received JSON to newUser
	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	// add new user
	users = append(users, newUser)
	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, newUser)
}

