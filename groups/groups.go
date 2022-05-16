package groups

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)


type Groups struct {
        Groups []Group `json:"groups"`
}

type Group struct {
        ID              string `json:"id"`
        Name	        string `json:"nickname"`
        Role            string `json:"role"`
}


// groups demo data for group struct
var groups = []Group{
	{ID: "1", Name: "superadmins"},
	{ID: "2", Name: "devs"},
	{ID: "3", Name: "ops"},
}

// GetGroups returns JSON serialized list of groups and their properties.
func GetGroups(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, groups)
}

// GetGroupByID returns group's properties, given sent ID exists in database.
func GetGroupByID(c *gin.Context) {
	id := c.Param("id")

	// loop over groups
	for _, a := range groups {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "group not found"})
}

// PostGroup enables one to add new group to demo groups model data.
func PostGroup(c *gin.Context) {
	var newGroup Group

	// bind received JSON to newGroup
	if err := c.BindJSON(&newGroup); err != nil {
		return
	}

	// add new group
	groups = append(groups, newGroup)
	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, newGroup)
}
