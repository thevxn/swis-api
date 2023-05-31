package roles

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var r sync.Map

// @Summary Get all roles
// @Description get roules complete list
// @Tags roles
// @Produce  json
// @Success 200 {object} roles.Roles
// @Router /roles [get]
// GetRoles returns JSON serialized list of roles and their properties.
func GetRoles(c *gin.Context) {
	var roles = make(map[string]Role)

	r.Range(func(rawKey, rawVal interface{}) bool {
		k, ok := rawKey.(string)
		v, ok := rawVal.(Role)

		if !ok {
			return false
		}

		roles[k] = v
		return true
	})

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "dumping roles",
		"roles":   roles,
	})
	return
}

// @Summary Get role by Name
// @Description get role by :id param
// @Tags roles
// @Produce  json
// @Success 200 {object} roles.Role
// @Router /roles/{name} [get]
func GetRoleByName(c *gin.Context) {
	var name string = c.Param("name")
	var role Role

	rawRole, ok := r.Load(name)
	role, ok = rawRole.(Role)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "role not found",
			"code":    http.StatusNotFound,
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "dumping requested role's details",
		"role":    role,
	})
	return
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
	var newRole *Role = &Role{}

	if err := c.BindJSON(newRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	if _, found := r.Load(newRole.Name); found {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "role already exists",
			"name":    newRole.Name,
		})
		return
	}

	r.Store(newRole.Name, newRole)

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "new role added",
		"role":    newRole,
	})
	return
}

// @Summary Update role by its Name
// @Description update role by its Name
// @Tags roles
// @Produce json
// @Param request body roles.Role.Name true "query params"
// @Success 200 {object} roles.Role
// @Router /roles/{name} [put]
func UpdateRoleByName(c *gin.Context) {
	var updatedRole *Role = &Role{}
	var name string = c.Param("name")

	if _, ok := r.Load(name); !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "role not found",
			"code":    http.StatusNotFound,
		})
		return
	}

	if err := c.BindJSON(updatedRole); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	r.Store(name, updatedRole)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "role updated",
		"role":    updatedRole,
	})
	return
}

// @Summary Delete role by its Name
// @Description delete role by its Name
// @Tags roles
// @Produce json
// @Param  id  path  string  true  "role Name"
// @Success 200 {object} roles.Role.Name
// @Router /roles/{name} [delete]
func DeleteRoleByName(c *gin.Context) {
	var name string = c.Param("name")

	r.Delete(name)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "role deleted by ID",
		"name":    name,
	})
	return
}

// @Summary Upload roles dump backup -- restores all roles
// @Description update roles JSON dump
// @Tags roles
// @Accept json
// @Produce json
// @Router /roles/restore [post]
// PostDumpRestore
func PostDumpRestore(c *gin.Context) {
	var importRoles = &Roles{}
	var role Role

	if err := c.BindJSON(importRoles); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	for _, role = range importRoles.Roles {
		r.Store(role.Name, role)
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "roles imported/restored, omitting output",
	})
	return
}
