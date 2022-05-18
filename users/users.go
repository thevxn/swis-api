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


func findUserByID(c *gin.Context) (index *int, u *User) {
	// loop over users
	for i, a := range users {
		if a.ID == c.Param("id") {
			//c.IndentedJSON(http.StatusOK, a)
			return &i, &a
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code": http.StatusNotFound,
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

// GetUserByID returns user's properties, given sent ID exists in database.
func GetUserByID(c *gin.Context) {
	//id := c.Param("id")

	if _, user := findUserByID(c); user != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	// add new user
	users = append(users, newUser)

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated,
		"message": "user added",
		"user": newUser,
	})
}

// PostUsersDumpRestore
func PostUsersDumpRestore(c *gin.Context) {
	var importUsers Users

	
	// bind received JSON to newUser
	if err := c.BindJSON(&importUsers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	// add new user
	users = importUsers.Users
	//users = append(users, newUser)

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated,
		"message": "users imported successfully",
	})
}

// PostUserSSHKey need "id" param
func PostUserSSHKey(c *gin.Context) {
	//var index *int, user *User = findUserByID(c)
	var index, user = findUserByID(c)

	// load SSH keys from POST request
	if err := c.BindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	// write changes to users array
	users[*index] = *user	
	c.IndentedJSON(http.StatusAccepted, *user)
}

