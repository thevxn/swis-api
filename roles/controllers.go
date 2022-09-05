package roles

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get all roles
// @Description get roules complete list
// @Tags roles
// @Produce  json
// @Success 200 {object} roles.Roles
// @Router /roles [get]
// GetGroups returns JSON serialized list of roles and their properties.
func GetRoles(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, roles)
}

// @Summary Get role by Name
// @Description get role by :id param
// @Tags roles
// @Produce  json
// @Success 200 {object} roles.Role
// @Router /role/{name} [get]
// GetGroupByID returns role's properties, given sent Name exists in database.
func GetRoleByName(c *gin.Context) {
	// loop over roles
	for _, r := range roles {
		if r.Name == c.Param("name") {
			c.IndentedJSON(http.StatusOK, r)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "role not found"})
}

// @Summary Add new role to roles array
// @Description add new role to roles array
// @Tags roles
// @Produce json
// @Param request body roles.Role true "query params"
// @Success 200 {object} roles.Role
// @Router /roles [post]
// PostGroup enables one to add new role to roles.
func PostRole(c *gin.Context) {
	var newRole Role

	// bind received JSON to newRole
	if err := c.BindJSON(&newRole); err != nil {
		return
	}

	// add new role
	roles = append(roles, newRole)
	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, newRole)
}

// @Summary Upload roles dump backup -- restores all roles
// @Description update roles JSON dump
// @Tags roles
// @Accept json
// @Produce json
// @Router /roles/restore [post]
// PostDumpRestore
func PostDumpRestore(c *gin.Context) {
	var importRoles Roles

	if err := c.BindJSON(&importRoles); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	// restore all roles
	roles = importRoles.Roles

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "roles imported successfully",
	})
}
