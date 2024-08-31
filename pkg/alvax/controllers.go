package alvax

import (
	"go.vxn.dev/swis/v5/pkg/core"

	"github.com/gin-gonic/gin"
)

var (
	Cache   *core.Cache
	pkgName string = "alvax"
)

var Package *core.Package = &core.Package{
	Name: pkgName,
	Cache: []**core.Cache{
		&Cache,
	},
	Routes:  Routes,
	Generic: true,
}

// GetConfigs function dumps the alvax cache contents.
// @Summary Get all alvax configs
// @Description get alvax config list
// @Tags alvax
// @Produce json
// @Success 200 {object} []alvax.ConfigRoot
// @Router /alvax [get]
func GetConfigs(ctx *gin.Context) {
	core.PrintAllRootItems(ctx, Cache, pkgName)
	return
}

// GetConfigByKey
// @Summary Get alvax config by key
// @Description get alvax config by key
// @Tags alvax
// @Produce json
// @Success 200 {object} alvax.ConfigRoot
// @Router /alvax/{key} [get]
func GetConfigByKey(ctx *gin.Context) {
	core.PrintItemByParam[ConfigRoot](ctx, Cache, pkgName, ConfigRoot{})
	return
}

// @Summary Add new alvax config
// @Description add new alvax config
// @Tags alvax
// @Produce json
// @Param request body alvax.ConfigRoot true "query params"
// @Success 200 {object} alvax.ConfigRoot
// @Router /alvax [post]
func PostNewConfig(ctx *gin.Context) {
	core.AddNewItem[ConfigRoot](ctx, Cache, pkgName, ConfigRoot{})
	return
}

// @Summary Update alvax config by its ID
// @Description update alvax config by its ID
// @Tags alvax
// @Produce json
// @Param request body alvax.ConfigRoot.Key true "query params"
// @Success 200 {object} alvax.ConfigRoot
// @Router /alvax/{key} [put]
func UpdateConfigByKey(ctx *gin.Context) {
	core.UpdateItemByParam[ConfigRoot](ctx, Cache, pkgName, ConfigRoot{})
	return
}

// @Summary Delete alvax config by its key
// @Description delete alvax config by its key
// @Tags alvax
// @Produce json
// @Param id path string true "alvax config key"
// @Success 200 {object} alvax.ConfigRoot.Key
// @Router /alvax/{key} [delete]
func DeleteConfigByKey(ctx *gin.Context) {
	core.DeleteItemByParam(ctx, Cache, pkgName)
	return
}

// PostDumpRestore
// @Summary Upload alvax configs dump -- restore configs
// @Description upload alvax config JSON dump and restore the data model
// @Tags alvax
// @Accept json
// @Produce json
// @Router /alvax/restore [post]
func PostDumpRestore(ctx *gin.Context) {
	core.BatchRestoreItems[ConfigRoot](ctx, Cache, pkgName)
	return
}

// @Summary List package model's field types
// @Description list package model's field types
// @Tags alvax
// @Accept json
// @Produce json
// @Router /alvax/types [get]
func ListTypes(ctx *gin.Context) {
	core.ParsePackageType(ctx, pkgName, ConfigRoot{})
	return
}
