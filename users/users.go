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
	ID       	string 	`json:"id"`
	Nickname 	string 	`json:"nickname"`
	Role     	string 	`json:"role"`
	TokenBase64	string 	`json:"token_base64"`
	SSHKeys	      []string  `json:"ssh_keys"`
	GPGKeys	      []string  `json:"gpg_keys"`
}

// users demo data for user struct
var users = []User{
	{ID: "1", Nickname: "sysadmin", Role: "admin"},
	{ID: "2", Nickname: "dev", Role: "developer"},
	{ID: "3", Nickname: "op", Role: "operator"},
}


func findUserByID(c *gin.Context) (u *User) {
	// loop over users
	for _, a := range users {
		if a.ID == c.Param("id") {
			//c.IndentedJSON(http.StatusOK, a)
			return &a
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
	return nil
}


// GetUsers returns JSON serialized list of users and their properties.
func GetUsers(c *gin.Context) {
	// serialize struct to JSON
	// net/http response code
	c.IndentedJSON(http.StatusOK, users)
}

// GetUserByID returns user's properties, given sent ID exists in database.
func GetUserByID(c *gin.Context) {
	//id := c.Param("id")

	if user := findUserByID(c); user != nil {
		// user found
		c.IndentedJSON(http.StatusOK, user)
	}

	//c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
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

// PostUserSSHKey need "id" param
func PostUserSSHKey(c *gin.Context) {
	var user *User = findUserByID(c)

	// load SSH keys from POST request
	if err := c.BindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	for _, a := range users {
		if a.ID == c.Param("id") {
			// save SSH keys to user
			a = *user
			c.IndentedJSON(http.StatusAccepted, user)
		}
	}

	c.IndentedJSON(http.StatusNotFound, user)
}

