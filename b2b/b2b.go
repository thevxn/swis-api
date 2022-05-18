package b2b

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)


type BusinessArray struct {
	BusinessArray []Business `json:"business_array"`
}

// Business structure
type Business struct {
	ID       		string 	`json:"id"`
	ICO			string  `json:"ico"`
	VAT			string  `json:"vat_id"`
	NameLabel 		string 	`json:"name_label"`
	Address		      []string  `json:"address"`
	Contact		     []Contact  `json:"contact"`
	CounterBankAccNo	string	`json:"account_no"`
	WebsiteURL		string  `json:"website_url"`
	Role     		string 	`json:"role"`
	TokenBase64		string 	`json:"token_base64"`
}

type Contact struct {
	Type	string 	`json:"type"`
	Content	string	`json:"content"`
}

// users demo data for user struct
var businessArray = []Business{
	{ID: "0", ICO: "1444968", VAT: "CZ9605144648", NameLabel: "Bc. Kryštof Šara", Address: []string{"Petrov 504, 696 65 Petrov"}, 
		Contact: []Contact{Type: "phone", Content: "+420 728 535 909"}},
	{ID: "1", ICO: "29321824", VAT: "CZ29321824", NameLabel: "Mathesio, s. r. o.", Address: []string{"Soukenická 558/3, Staré Brno, 602 00 Brno"}},
}


func findBusinessByID(c *gin.Context) (b *Business) {
	// loop over businesses
	for _, b := range businessArray {
		if b.ID == c.Param("id") {
			//c.IndentedJSON(http.StatusOK, a)
			return &b
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code": http.StatusNotFound,
		"message": "business not found",
	})
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

