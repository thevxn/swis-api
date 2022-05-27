package infra

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Infrastructure struct {
	//Users		[]User
	//Groups	[]Group
	Hosts		[]Host
	Networks	[]Network
}


type Network struct {
	Name		string	`json:"network_name"`
	Address		string	`json:"network_address"`
	CIDRBlock	string	`json:"network_cidr_block"`	
}

/*
 * IP address alocation (CIDR) count + possible supernetwork "expandability"
 * https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing#IPv4_CIDR_blocks
 *
 * /24 -- 256-2
 * /25 -- 128-2
 * /26 -- 64-2 (bigger lan networks)
 * /27 -- 32-2 (small lan networks)
 *
 */

var networks = []Network{
	{Name: "squabbit virbr0 network", Address: "10.4.5.0", CIDRBlock: "27"},
	{Name: "-vacant virbr32 network", Address: "10.4.5.32", CIDRBlock: "27"},
	{Name: "-vacant virbr32 network", Address: "10.4.5.64", CIDRBlock: "27"},
	{Name: "-vacant network", Address: "10.4.5.96", CIDRBlock: "27"},
	{Name: "-vacant network", Address: "10.4.5.128", CIDRBlock: "27"},
	{Name: "VPN external dish network", Address: "10.4.5.160", CIDRBlock: "27"},
	{Name: "VPN private intranet", Address: "10.4.5.192", CIDRBlock: "27"},
	{Name: "VPN client route-all-traffic network", Address: "10.4.5.224", CIDRBlock: "27"},
}

type Hosts struct {
	Hosts []Host `json:"hosts"`
}

type Host struct {
	ID       	string 	`json:"id"`
	Hostname 	string 	`json:"nickname"`
	Domain		string	`json:"domain"`
	Role     	string 	`json:"role"`
	IPAddress     []string 	`json:"ip_address"`
	Facter	      []string	`json:"facter"`
	VMs	      []Host	`json:"virtual_machines"`
}

// Hyper struct to model hypervisor machine
type Hyper struct {
	Host
}

// Virtual struct to model virtual machine
type Virtual struct {
	Host
}

var squabbitVMs []Host{
	{Hostname: "stokrle", Domain: "savla.su", Role: "build,deploy", IPAddress: []string{"10.4.5.55/25"}},
}

// demo Hosts data
var hosts Hosts{
	{Hostname: "squabbit", Domain: "savla.su", Role: "hypervisor", IPAddress: []string{"10.4.5.1/25", "10.4.5.129/25"}, VMs},
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

