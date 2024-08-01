package business

import (
	"github.com/gin-gonic/gin"

	"go.savla.dev/swis/v5/pkg/core"
)

var (
	Cache   *core.Cache
	pkgName string = "business"
)

var Package *core.Package = &core.Package{
	Name: pkgName,
	Cache: []**core.Cache{
		&Cache,
	},
	Routes:  Routes,
	Generic: true,
}

// @Summary Get all business entities
// @Description get business entities list
// @Tags business
// @Produce  json
// @Success 200 {object} []business.Business
// @Router /business [get]
func GetBusinessEntities(ctx *gin.Context) {
	core.PrintAllRootItems(ctx, Cache, pkgName)
	return
}

// @Summary Get business entity by its key
// @Description get business by key param
// @Tags business
// @Produce  json
// @Success 200 {object} business.Business
// @Router /business/{key} [get]
func GetBusinessByKey(ctx *gin.Context) {
	core.PrintItemByParam[Business](ctx, Cache, pkgName, Business{})
	return
}

// @Summary Add new business entity
// @Description add new business entity
// @Tags business
// @Produce json
// @Param request body business.Business true "query params"
// @Success 200 {object} business.Business
// @Router /business/{key} [post]
func PostBusinessByKey(ctx *gin.Context) {
	core.AddNewItemByParam[Business](ctx, Cache, pkgName, Business{})
	return
}

// @Summary Update business entity by its key
// @Description update business entity by its key
// @Tags business
// @Produce json
// @Param request body business.Business.ID true "query params"
// @Success 200 {object} business.Business
// @Router /business/{key} [put]
func UpdateBusinessByKey(ctx *gin.Context) {
	core.UpdateItemByParam[Business](ctx, Cache, pkgName, Business{})
	return
}

// @Summary Delete business by its key
// @Description delete business by its key
// @Tags business
// @Produce json
// @Param  id  path  string  true  "business ID"
// @Success 200 {object} business.Business.ID
// @Router /business/{key} [delete]
func DeleteBusinessByKey(ctx *gin.Context) {
	core.DeleteItemByParam(ctx, Cache, pkgName)
	return
}

// PostDumpRestore
// @Summary Upload business dump backup -- restores all business entities
// @Description upload business JSON dump
// @Tags business
// @Accept json
// @Produce json
// @Router /business/restore [post]
func PostDumpRestore(ctx *gin.Context) {
	core.BatchRestoreItems[Business](ctx, Cache, pkgName)
	return
}
