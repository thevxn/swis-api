package business

import (
	"github.com/gin-gonic/gin"

	"go.savla.dev/swis/v5/config"
)

var (
	Cache   *config.Cache
	pkgName string = "business"
)

// @Summary Get all business entities
// @Description get business entities list
// @Tags business
// @Produce  json
// @Success 200 {object} []business.Business
// @Router /business [get]
func GetBusinessEntities(ctx *gin.Context) {
	config.PrintAllRootItems(ctx, Cache, pkgName)
	return
}

// @Summary Get business entity by its key
// @Description get business by key param
// @Tags business
// @Produce  json
// @Success 200 {object} business.Business
// @Router /business/{key} [get]
func GetBusinessByKey(ctx *gin.Context) {
	config.PrintItemByParam(ctx, Cache, pkgName, Business{})
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
	config.AddNewItemByParam(ctx, Cache, pkgName, Business{})
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
	config.UpdateItemByParam(ctx, Cache, pkgName, Business{})
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
	config.DeleteItemByParam(ctx, Cache, pkgName)
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
	config.BatchRestoreItems(ctx, Cache, pkgName, Business{})
	return
}
