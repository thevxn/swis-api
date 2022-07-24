package users

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func findUserByName(c *gin.Context) (index *int, u *User) {
	// loop over users
	for i, a := range users {
		if a.Name == c.Param("name") {
			//c.IndentedJSON(http.StatusOK, a)
			return &i, &a
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "user not found",
	})
	return nil, nil
}

// GetUsers returns JSON serialized list of users and their properties.
func GetUsers(c *gin.Context) {
	// serialize struct to JSON
	// net/http response code
	c.IndentedJSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// GetUserByName returns user's properties, given sent name exists in database.
func GetUserByName(c *gin.Context) {
	if _, user := findUserByName(c); user != nil {
		// user found
		c.IndentedJSON(http.StatusOK, user)
	}

	//c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

// PostNewUser enables one to add new user to users model.
func PostNewUser(c *gin.Context) {
	var newUser User

	// bind received JSON to newUser
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	// add new user
	users = append(users, newUser)

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "user added",
		"user":    newUser,
	})
}

// PostDumpRestore
func PostDumpRestore(c *gin.Context) {
	var importUsers Users

	// bind received JSON to newUser
	if err := c.BindJSON(&importUsers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	// add new user
	users = importUsers.Users
	//users = append(users, newUser)

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "users imported successfully",
	})
}

// PostUsersSSHKeys method adds (rewrites) SSH key array by user.Name
func PostUsersSSHKeys(c *gin.Context) {
	var index, user = findUserByName(c)

	// load SSH keys from POST request
	if err := c.BindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	// write changes to users array
	users[*index] = *user
	c.IndentedJSON(http.StatusAccepted, *user)
}

// GetUsersSSHKeysRaw
func GetUsersSSHKeysRaw(c *gin.Context) {
	var _, user = findUserByName(c)

	if user != nil {
		var responseBody = strings.Join(user.SSHKeys, "\n")
		c.String(http.StatusOK, responseBody)
	}
	return
}
