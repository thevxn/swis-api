// standalone program 
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// in-memory Model -- declaration
type user struct {
	ID    	  string `json:"id"`
	Nickname  string `json:"nickname"`
	Role	  string `json:"role"`
}

// data
var users = []user{
	{ID: "1", Nickname: "sysadmin", Role: "admin"},
	{ID: "2", Nickname: "dev", Role: "developer"},
	{ID: "3", Nickname: "op", Role: "operator"},
}

/*
 * HANDLERS
 */
// reguest details stored in gin.Context
func getUsers(c *gin.Context) {
	// serialize struct to JSON
	// net/http response code 
	c.IndentedJSON(http.StatusOK, users)
}

func getUserByID(c *gin.Context) {
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

func postUser(c *gin.Context) {
	var newUser user

	// bind received JSON to newUser
	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	// add new user
	users = append(users, newUser)
	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, newUser)
}

func main() {
	// name of the function, not its result
	router := gin.Default() 
	router.GET("/users", getUsers)
	router.GET("/users/:id", getUserByID)
	router.POST("/users", postUser)

	// attach router to http.Server and start it
	router.Run("0.0.0.0:8080")
}
