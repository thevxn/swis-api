package users

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var u sync.Map

func FindUserByToken(token string) *User {
	var users = make(map[string]User)
	var user *User = nil

	u.Range(func(rawKey, rawVal interface{}) bool {
		// very insecure assert
		k, ok := rawKey.(string)
		v, ok := rawVal.(User)

		if !ok {
			return false
		}

		// each token should be unique, to be generated from User.Name and other attributes + pepper
		if v.TokenHMAC == token && v.Active {
			users[k] = v

			// if one token is shared between users, disallow both (this to be more discussed, please)
			if len(users) > 1 {
				return false
			}

			user = &v
		}

		return true
	})

	return user
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
	var users = make(map[string]User)

	u.Range(func(rawKey, rawVal interface{}) bool {
		// very insecure assert
		k, ok := rawKey.(string)
		v, ok := rawVal.(User)

		if !ok {
			return false
		}

		users[k] = v
		return true
	})

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
	var name string = c.Param("name")
	var user User

	userRaw, ok := u.Load(name)
	user, ok = userRaw.(User)

	if !ok {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"message": "user not found or the user object cannot be loaded",
			"code":    http.StatusConflict,
			"name":    name,
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "ok, dumping user info",
		"code":    http.StatusOK,
		"user":    user,
	})
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
	var newUser = &User{}

	// bind received JSON to newUser
	if err := c.BindJSON(newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	if _, found := u.Load(newUser.Name); found {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"message": "user already exists",
			"code":    http.StatusConflict,
			"name":    newUser.Name,
		})
		return
	}

	u.Store(newUser.Name, newUser)

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
	var importUsers = &Users{}

	// bind received JSON to newUser
	if err := c.BindJSON(importUsers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}
	for _, user := range importUsers.Users {
		u.Store(user.Name, user)
	}

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "users imported successfully",
	})
}

// (PUT /users/{name}/active)
// @Summary Toggle active boolean for {user}
// @Description toggle active boolean for {user}
// @Tags users
// @Produce json
// @Param  id  path  string  true  "username"
// @Success 200 {object} users.User
// @Router /users/{name}/active [put]
func ActiveToggleUserByName(c *gin.Context) {
	var user User
	var userName string = c.Param("name")

	rawUser, ok := u.Load(userName)
	user, ok = rawUser.(User)

	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "user not found",
			"code":    http.StatusNotFound,
			"name":    userName,
		})
		return
	}

	// inverse the Active field value
	user.Active = !user.Active

	u.Store(userName, user)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "user active toggle pressed!",
		"user":    user,
	})
	return
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
	var user User
	var userName string = c.Param("name")

	rawUser, ok := u.Load(userName)
	user, ok = rawUser.(User)

	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "user not found",
			"code":    http.StatusNotFound,
			"name":    userName,
		})
		return
	}

	// to be reimplemented later
	var sshKeys []string

	// load SSH keys from POST request
	if err := c.BindJSON(&sshKeys); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	user.SSHKeys = sshKeys
	u.Store(userName, user)

	c.IndentedJSON(http.StatusAccepted, gin.H{
		"code":    http.StatusAccepted,
		"message": "ssh keys for user imported",
		"name":    userName,
	})
	return
}

// @Summary Get User's SSH keys in plain text
// @Description fetch :user ssh key array output in plain text
// @Tags users
// @Produce json
// @Success 200 {object} users.User
// @Router /users/{name}/keys/ssh [get]
// GetUsersSSHKeysRaw
func GetUsersSSHKeysRaw(c *gin.Context) {
	var user User
	var userName string = c.Param("name")

	rawUser, ok := u.Load(userName)
	user, ok = rawUser.(User)

	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "user not found",
			"code":    http.StatusNotFound,
			"name":    userName,
		})
		return
	}

	// return SSH keys as plaintext
	var responseBody = strings.Join(user.SSHKeys, "\n")
	c.String(http.StatusOK, responseBody)

	return
}

// @Summary Delete user by Name
// @Description delete user by Name
// @Tags users
// @Produce json
// @Param  id  path  string  true  "user Name"
// @Success 200 {object} users.User.Name
// @Router /users/{name} [delete]
func DeleteUserByName(c *gin.Context) {
	var name string = c.Param("name")

	u.Delete(name)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "user deleted by Name",
		"name":    name,
	})
	return
}

// @Summary Update user by Name
// @Description update user by Name
// @Tags users
// @Produce json
// @Param request body users.User.Name true "query params"
// @Success 200 {object} users.User
// @Router /users/{name} [put]
func UpdateUserByName(c *gin.Context) {
	var user = &User{}
	var name string = c.Param("name")

	if _, ok := u.Load(name); !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "project not found",
			"code":    http.StatusNotFound,
		})
		return
	}

	if err := c.BindJSON(user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	u.Store(name, user)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "user updated",
		"user":    user,
	})
	return
}
