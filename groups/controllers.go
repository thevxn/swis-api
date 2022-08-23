package groups

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get all groups
// @Description get groups complete list
// @Tags groups
// @Produce  json
// @Success 200 {object} groups.Groups
// @Router /groups [get]
// GetGroups returns JSON serialized list of groups and their properties.
func GetGroups(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, groups)
}

// @Summary Get group by ID
// @Description get group by :id param
// @Tags groups
// @Produce  json
// @Success 200 {object} groups.Group
// @Router /groups/{id} [get]
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

// @Summary Add new group to groups array
// @Description add new group to groups array
// @Tags groups
// @Produce json
// @Param request body groups.Group true "query params"
// @Success 200 {object} groups.Group
// @Router /group [post]
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
