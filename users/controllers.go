package users

import (
	"net/http"
	"strings"

	"go.savla.dev/swis/v5/config"

	"github.com/gin-gonic/gin"
)

var Cache *config.Cache

func FindUserByToken(token string) *User {
	rawUsers, _ := Cache.GetAll()

	for _, rawUser := range rawUsers {
		user, ok := rawUser.(User)
		if !ok {
			return nil
		}

		if user.TokenHash == token && user.Active {
			return &user
		}
	}

	return nil
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
	users, count := Cache.GetAll()

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"count":   count,
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
	var userName string = c.Param("name")
	var user User

	rawUser, ok := Cache.Get(userName)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "user not found",
			"name":    userName,
		})
		return
	}

	user, ok = rawUser.(User)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "cannot assert data type, database internal error",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, dumping user's info",
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
	var newUser User

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot bind input JSON stream",
		})
		return
	}

	// TODO: implement LoadOrStore() method
	if _, found := Cache.Get(newUser.Name); found {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "user already exists",
			"name":    newUser.Name,
		})
		return
	}

	if saved := Cache.Set(newUser.Name, newUser); !saved {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "user couldn't be saved to database",
		})
		return
	}

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
	var user User
	var counter int = 0

	// bind received JSON to newUser
	if err := c.BindJSON(importUsers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot bind input JSON stream",
		})
		return
	}

	for _, user = range importUsers.Users {
		Cache.Set(user.Name, user)
		counter++
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "users imported successfully",
		"count":   counter,
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

	rawUser, ok := Cache.Get(userName)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "user not found",
			"code":    http.StatusNotFound,
			"name":    userName,
		})
		return
	}

	user, ok = rawUser.(User)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "cannot assert data type, database internal error",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// inverse the Active field value
	user.Active = !user.Active

	if saved := Cache.Set(userName, user); !saved {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "user couldn't be saved to database",
		})
		return
	}

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

	rawUser, ok := Cache.Get(userName)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "user not found",
			"code":    http.StatusNotFound,
			"name":    userName,
		})
		return
	}

	user, ok = rawUser.(User)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "cannot assert data type, database internal error",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// to be reimplemented later
	var sshKeys []string

	// load SSH keys from POST request
	if err := c.BindJSON(&sshKeys); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot bind input JSON stream",
		})
		return
	}

	user.SSHKeys = sshKeys

	if saved := Cache.Set(userName, user); !saved {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "user couldn't be saved to database",
		})
		return
	}

	c.IndentedJSON(http.StatusAccepted, gin.H{
		"code":    http.StatusAccepted,
		"message": "ssh keys for user (re)imported",
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

	rawUser, ok := Cache.Get(userName)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "user not found",
			"code":    http.StatusNotFound,
			"name":    userName,
		})
		return
	}

	user, ok = rawUser.(User)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "cannot assert data type, database internal error",
			"code":    http.StatusInternalServerError,
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

	if _, found := Cache.Get(name); !found {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "user not found",
			"code":    http.StatusNotFound,
		})
		return
	}

	if deleted := Cache.Delete(name); !deleted {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "user couldn't be deleted from database",
		})
		return
	}

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
	var name string = c.Param("name")
	var updatedUser User

	if _, found := Cache.Get(name); !found {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "user not found",
		})
		return
	}

	if err := c.BindJSON(updatedUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot bind input JSON stream",
		})
		return
	}

	if saved := Cache.Set(name, updatedUser); !saved {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "user couldn't be saved to database",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "user updated",
		"user":    updatedUser,
	})
	return
}
