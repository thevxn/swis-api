package depots

import (
	"net/http"
	//"sort"
	"strconv"

	"go.savla.dev/swis/v5/pkg/core"

	"github.com/gin-gonic/gin"
)

var (
	Cache   *core.Cache
	pkgName string = "depots"
)

// GetSocketList GET method
// @Summary Get all depots and their users/owners
// @Description get depot complete list
// @Tags depots
// @Produce json
// @Success 200 {object} []depots.DepotItem
// @Router /depots [get]
func GetDepotItems(ctx *gin.Context) {
	core.PrintAllRootItems(ctx, Cache, pkgName)
	return
}

// @Summary Add new depot item
// @Description add new depot item
// @Tags depots
// @Produce json
// @Param request body depots.DepotItem true "query params"
// @Success 200 {object} depots.DepotItem
// @Router /depots/{key} [post]
func PostNewDepotItemByKey(ctx *gin.Context) {
	core.AddNewItemByParam(ctx, Cache, pkgName, DepotItem{})
	return
}

// @Summary Update depot item by its key
// @Description update depot by its key
// @Tags depots
// @Produce json
// @Param  id  path  string  true  "depot key"
// @Success 200 {object} depots.DepotItem
// @Router /depots/{key} [put]
func UpdateDepotItemByKey(ctx *gin.Context) {
	core.UpdateItemByParam(ctx, Cache, pkgName, DepotItem{})
	return
}

// @Summary Delete depot item by its key
// @Description delete depot item by its key
// @Tags depots
// @Produce json
// @Param  id  path  string  true  "depot key"
// @Success 200 {object} depots.DepotItem
// @Router /depots/{key} [delete]
func DeleteDepotItemByKey(ctx *gin.Context) {
	core.DeleteItemByParam(ctx, Cache, pkgName)
	return
}

// @Summary Upload depot dump backup -- restores all depot items
// @Description upload depots JSON dump
// @Tags depots
// @Accept json
// @Produce json
// @Router /depots/restore [post]
func PostDumpRestore(ctx *gin.Context) {
	core.BatchRestoreItems(ctx, Cache, pkgName, DepotItem{})
	return
}

// @Summary Get depot item list by Owner
// @Description get depot item list by :owner param
// @Tags depots
// @Produce json
// @Success 200 {object} []depots.DepotItem
// @Router /depots/{owner} [get]
func GetDepotItemsByOwner(ctx *gin.Context) {
	var owner string = ctx.Param("owner")
	var exportedItemsMap = make(map[int]DepotItem)

	rawItemsMap, _ := Cache.GetAll()

	for rawKey, rawItem := range rawItemsMap {
		key, err := strconv.Atoi(rawKey)
		if err != nil {
			continue
		}

		item, ok := rawItem.(DepotItem)
		if !ok {
			continue
		}

		if item.Owner == owner {
			// TODO: reimplement this
			// order items ASC alphabetically
			// https://pkg.go.dev/sort#Slice
			//sort.Slice(depot.DepotItems, func(i, j int) bool {
			//	return (depot.DepotItems[i].Description < depot.DepotItems[j].Description)
			//})

			exportedItemsMap[key] = item
		}
	}

	if len(exportedItemsMap) > 0 {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"items":   exportedItemsMap,
			"message": "ok, dumping items for owner",
			"owner":   owner,
			"package": pkgName,
		})
		return
	}

	ctx.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "no item found for such owner",
		"owner":   owner,
		"package": pkgName,
	})

	return
}
