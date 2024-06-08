package roles

import (
	"go.savla.dev/swis/v5/pkg/core"

	"github.com/gin-gonic/gin"
)

var (
	Cache   *core.Cache
	pkgName string = "roles"
)

var Package *core.Package = &core.Package{
	Name: pkgName,
	Cache: []**core.Cache{
		&Cache,
	},
	Routes: Routes,
}

// GetRoles returns JSON serialized list of roles and their properties.
// @Summary Get all roles
// @Description get roules complete list
// @Tags roles
// @Produce  json
// @Success 200 {object} []roles.Role
// @Router /roles [get]
func GetRoles(ctx *gin.Context) {
	core.PrintAllRootItems(ctx, Cache, pkgName)
	return
}

// @Summary Get role by Key
// @Description get role by :key param
// @Tags roles
// @Produce  json
// @Success 200 {object} roles.Role
// @Router /roles/{key} [get]
func GetRoleByKey(ctx *gin.Context) {
	core.PrintItemByParam[Role](ctx, Cache, pkgName)
	return
}

// PostNewRoleByKey enables one to add new role to roles.
// @Summary Add new role to roles array
// @Description add new role to roles array
// @Tags roles
// @Produce json
// @Param request body roles.Role true "query params"
// @Success 200 {object} roles.Role
// @Router /roles/{key} [post]
func PostNewRoleByKey(ctx *gin.Context) {
	core.AddNewItemByParam[Role](ctx, Cache, pkgName)
	return
}

// @Summary Update role by its Key
// @Description update role by its Key
// @Tags roles
// @Produce json
// @Param request body roles.Role.Name true "query params"
// @Success 200 {object} roles.Role
// @Router /roles/{key} [put]
func UpdateRoleByKey(ctx *gin.Context) {
	core.UpdateItemByParam[Role](ctx, Cache, pkgName)
	return
}

// @Summary Delete role by its Key
// @Description delete role by its Key
// @Tags roles
// @Produce json
// @Param id path string true "role Key"
// @Success 200 {object} roles.Role.Name
// @Router /roles/{key} [delete]
func DeleteRoleByKey(ctx *gin.Context) {
	core.DeleteItemByParam(ctx, Cache, pkgName)
	return
}

// PostDumpRestore
// @Summary Upload roles dump backup -- restores all roles
// @Description update roles JSON dump
// @Tags roles
// @Accept json
// @Produce json
// @Router /roles/restore [post]
func PostDumpRestore(ctx *gin.Context) {
	core.BatchRestoreItems[Role](ctx, Cache, pkgName)
	return
}
