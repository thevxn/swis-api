package users

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func FindUserByToken(token string) *User {
	// Loop over all loaded users.
	for _, u := range users {
		if u.TokenHMAC == token && u.Active {
			return &u
		}
	}
	return nil
}

// findUserByName is a private, hellper function for users array struct browsing.
// *gin.Context should contain the 'name' parameter (extracted from HTTP path string).
func findUserByName(c *gin.Context) (index *int, u *User) {
	// Loop over all loaded users.
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

// @Summary Get all users
// @Description get users complete list
// @Tags users
// @Produce  json
// @Success 200 {object} users.Users
// @Router /users [get]
// GetSocketList GET method
// GetUsers returns JSON serialized list of users and their properties.
func GetUsers(c *gin.Context) {
	// serialize struct to JSON
	// net/http response code
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, listing users",
		"users":   users,
	})
}

// @Summary Get user by Name
// @Description get user by their :name param
// @Tags users
// @Produce  json
// @Success 200 {object} users.User
// @Router /users/{name} [get]
// GetUserByName returns user's properties, given sent name exists in database.
func GetUserByName(c *gin.Context) {
	if _, user := findUserByName(c); user != nil {
		// user found
		c.IndentedJSON(http.StatusOK, user)
	}

	//c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

// @Summary Add new user to users array
// @Description add new user to users array
// @Tags users
// @Produce json
// @Param request body users.User true "query params"
// @Success 200 {object} users.User
// @Router /users [post]
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

// @Summary Upload users dump backup -- restores all users
// @Description update users JSON dump
// @Tags users
// @Accept json
// @Produce json
// @Router /users/restore [post]
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

// @Summary Add SSH public keys to User
// @Description add new SSH keys to :user param
// @Tags users
// @Produce json
// @Param request body string true "query params"
// @Success 200 {object} users.User
// @Router /users/{name}/keys/ssh [post]
// PostUsersSSHKeys method adds (rewrites) SSH key array by user.Name
func PostUsersSSHKeys(c *gin.Context) {
	// shoud be safe, findeUserByName should not return a nil
	var index, user = findUserByName(c)
	if user == nil {
		return
	}

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
	c.IndentedJSON(http.StatusAccepted, gin.H{
		"code":    http.StatusAccepted,
		"message": "ssh keys for user imported",
		"user":    *user,
	})
}

// @Summary Get User's SSH keys in plain text
// @Description fetch :user ssh key array output in plain text
// @Tags users
// @Produce json
// @Success 200 {object} users.User
// @Router /users/{name}/keys/ssh [get]
// GetUsersSSHKeysRaw
func GetUsersSSHKeysRaw(c *gin.Context) {
	var _, user = findUserByName(c)

	if user != nil {
		var responseBody = strings.Join(user.SSHKeys, "\n")
		c.String(http.StatusOK, responseBody)
	}
	return
}
