package roles

import (
	"go.vxn.dev/swis/v5/pkg/core"

	"github.com/gin-gonic/gin"
)

var (
	Cache *core.Cache

	caches = []**core.Cache{
		&Cache,
	}
	pkgName string = "roles"
)

var Package *core.Package = &core.Package{
	Name:    pkgName,
	Cache:   caches,
	Routes:  Routes,
	Generic: true,
}

var restorePackage = &core.RestorePackage{
	Name:             pkgName,
	Cache:            caches,
	CacheNames:       []string{"Cache"},
	Subpackages:      []string{},
	SubpackageModels: map[string]any{},
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
	core.PrintItemByParam[Role](ctx, Cache, pkgName, Role{})
	return
}

// PostNewRoleByKey enables one to add new role to roles.
// @Summary Add new role to roles array
// @Description add new role to roles array
// @Tags roles
// @Produce json
// @Param request body roles.Role true "query params"
// @Success 200 {object} roles.Role
// @Router /roles [post]
func PostNewRole(ctx *gin.Context) {
	core.AddNewItem[Role](ctx, Cache, pkgName, Role{})
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
	core.UpdateItemByParam[Role](ctx, Cache, pkgName, Role{})
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
	core.BatchRestoreItems[Role](ctx, restorePackage)
	return
}

// @Summary List package model's field types
// @Description list package model's field types
// @Tags roles
// @Accept json
// @Produce json
// @Router /roles/types [get]
func ListTypes(ctx *gin.Context) {
	core.ParsePackageType(ctx, pkgName, Role{})
	return
}
