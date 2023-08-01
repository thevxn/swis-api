package roles

import (
	"net/http"

	"go.savla.dev/swis/v5/config"

	"github.com/gin-gonic/gin"
)

var Cache *config.Cache

// @Summary Get all roles
// @Description get roules complete list
// @Tags roles
// @Produce  json
// @Success 200 {object} roles.Roles
// @Router /roles [get]
// GetRoles returns JSON serialized list of roles and their properties.
func GetRoles(c *gin.Context) {
	var roles = Cache.GetAll()

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, dumping all roles",
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

	rawRole, ok := Cache.Get(name)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "role not found",
			"code":    http.StatusNotFound,
		})
		return
	}

	role, ok = rawRole.(Role)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "cannot assert data type, database internal error",
			"code":    http.StatusInternalServerError,
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
// PostNewRole enables one to add new role to roles.
func PostNewRole(c *gin.Context) {
	var newRole Role

	if err := c.BindJSON(newRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot bind input JSON stream",
		})
		return
	}

	if _, found := Cache.Get(newRole.Name); found {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "role already exists",
			"name":    newRole.Name,
		})
		return
	}

	if saved := Cache.Set(newRole.Name, newRole); !saved {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "role couldn't be saved to database",
		})
		return
	}

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
	var name string = c.Param("name")
	var updatedRole Role

	if _, found := Cache.Get(name); !found {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "role not found",
			"code":    http.StatusNotFound,
		})
		return
	}

	if err := c.BindJSON(updatedRole); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot bind input JSON stream",
		})
		return
	}

	if saved := Cache.Set(name, updatedRole); !saved {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "role couldn't be saved to database",
		})
		return
	}

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
// @Param id path string true "role Name"
// @Success 200 {object} roles.Role.Name
// @Router /roles/{name} [delete]
func DeleteRoleByName(c *gin.Context) {
	var name string = c.Param("name")

	if _, found := Cache.Get(name); !found {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "role not found",
			"code":    http.StatusNotFound,
		})
		return
	}

	if deleted := Cache.Delete(name); !deleted {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "role couldn't be deleted from database",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "role deleted by Name",
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
	var counter int = 0

	if err := c.BindJSON(importRoles); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot bind input JSON stream",
		})
		return
	}

	for _, role = range importRoles.Roles {
		Cache.Set(role.Name, role)
		counter++
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "roles imported/restored, omitting output",
		"count":   counter,
	})
	return
}
