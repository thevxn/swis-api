package users

import (
	"net/http"
	"strings"

	"go.savla.dev/swis/v5/pkg/core"

	"github.com/gin-gonic/gin"
)

var (
	Cache   *core.Cache
	pkgName string = "users"
)

var Package *core.Package = &core.Package{
	Name: pkgName,
	Cache: []**core.Cache{
		&Cache,
	},
	Routes: Routes,
}

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

// GetUsers returns JSON serialized list of users and their properties.
// @Summary Get all users
// @Description get users complete list
// @Tags users
// @Produce  json
// @Success 200 {object} []users.User
// @Router /users [get]
// GetSocketList GET method
func GetUsers(ctx *gin.Context) {
	core.PrintAllRootItems(ctx, Cache, pkgName)
	return
}

// GetUserByName returns user's properties, given sent name exists in database.
// @Summary Get user by Key
// @Description get user by their :key param
// @Tags users
// @Produce  json
// @Success 200 {object} users.User
// @Router /users/{key} [get]
func GetUserByKey(ctx *gin.Context) {
	core.PrintItemByParam[User](ctx, Cache, pkgName)
	return
}

// PostNewUserByKey enables one to add new user to users model.
// @Summary Add new user to users array
// @Description add new user to users array
// @Tags users
// @Produce json
// @Param request body users.User true "query params"
// @Success 200 {object} users.User
// @Router /users/{key} [post]
func PostNewUserByKey(ctx *gin.Context) {
	core.AddNewItemByParam[User](ctx, Cache, pkgName)
	return
}

// @Summary Update user by Key
// @Description update user by Key
// @Tags users
// @Produce json
// @Param request body users.User.Name true "query params"
// @Success 200 {object} users.User
// @Router /users/{key} [put]
func UpdateUserByKey(ctx *gin.Context) {
	core.UpdateItemByParam[User](ctx, Cache, pkgName)
	return
}

// @Summary Delete user by Key
// @Description delete user by Key
// @Tags users
// @Produce json
// @Param  id  path  string  true  "user Name"
// @Success 200 {object} users.User.Name
// @Router /users/{key} [delete]
func DeleteUserByKey(ctx *gin.Context) {
	core.DeleteItemByParam(ctx, Cache, pkgName)
	return
}

// @Summary Upload users dump backup -- restores all users
// @Description update users JSON dump
// @Tags users
// @Accept json
// @Produce json
// @Router /users/restore [post]
// PostDumpRestore
func PostDumpRestore(ctx *gin.Context) {
	core.BatchRestoreItems[User](ctx, Cache, pkgName)
	return
}

// (PUT /users/{name}/active)
// @Summary Toggle active boolean for {user}
// @Description toggle active boolean for {user}
// @Tags users
// @Produce json
// @Param  id  path  string  true  "user name"
// @Success 200 {object} users.User
// @Router /users/{key}/active [put]
func ActiveToggleUserByKey(c *gin.Context) {
	var user User
	var userName string = c.Param("key")

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
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
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

// PostUsersSSHKeys method adds (rewrites) SSH key array by user.Name.
// @Summary Add SSH public keys to User
// @Description add new SSH keys to :user param
// @Tags users
// @Produce json
// @Param request body string true "query params"
// @Success 200 {object} users.User
// @Router /users/{key}/keys/ssh [post]
func PostUsersSSHKeys(c *gin.Context) {
	var user User
	var userName string = c.Param("key")

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
	var sshKeys struct {
		Keys []string `json:"keys"`
	}

	// load SSH keys from POST request
	if err := c.BindJSON(&sshKeys); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot bind input JSON stream",
		})
		return
	}

	user.SSHKeys = sshKeys.Keys

	if saved := Cache.Set(userName, user); !saved {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "user couldn't be saved to database",
		})
		return
	}

	c.IndentedJSON(http.StatusAccepted, gin.H{
		"code":    http.StatusAccepted,
		"message": "ssh keys for user (re)imported",
		"name":    userName,
		"user":    user,
	})
	return
}

// GetUsersSSHKeysRaw
// @Summary Get User's SSH keys in plain text
// @Description fetch :user ssh key array output in plain text
// @Tags users
// @Produce json
// @Success 200 {object} users.User
// @Router /users/{key}/keys/ssh [get]
func GetUsersSSHKeysRaw(c *gin.Context) {
	var user User
	var userName string = c.Param("key")

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
