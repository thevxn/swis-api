package links

import (
	"net/http"

	"go.savla.dev/swis/v5/config"

	"github.com/gin-gonic/gin"
)

var (
	Cache   *config.Cache
	pkgName string = "links"
)

// GetLinks returns JSON serialized list of links and their properties.
// @Summary Get all links
// @Description get links complete list
// @Tags links
// @Produce json
// @Success 200 {object} links.Link
// @Router /links [get]
func GetLinks(ctx *gin.Context) {
	config.PrintAllRootItems(ctx, Cache, pkgName)
	return
}

// GetLinkByKey returns link's properties, given sent hash exists in database.
// @Summary Get link by :hash
// @Description get link by its :hash param
// @Tags links
// @Produce json
// @Success 200 {object} links.Link
// @Router /links/{key} [get]
func GetLinkByKey(ctx *gin.Context) {
	config.PrintItemByParam(ctx, Cache, pkgName, Link{})
	return
}

// PostNewLinkByKey enables one to add new link to links model.
// @Summary Add new link to links
// @Description add new link to links array
// @Tags links
// @Produce json
// @Param request body links.Link true "query params"
// @Success 200 {object} links.Link
// @Router /links/{key} [post]
func PostNewLinkByKey(ctx *gin.Context) {
	config.AddNewItemByParam(ctx, Cache, pkgName, Link{})
	return
}

// @Summary Upload links dump backup -- restore all links
// @Description update links JSON dump
// @Tags links
// @Accept json
// @Produce json
// @Router /links/restore [post]
func PostDumpRestore(ctx *gin.Context) {
	config.BatchRestoreItems(ctx, Cache, pkgName, Link{})
	return
}

// @Summary Update link by its Key
// @Description update link by its Key
// @Tags links
// @Produce json
// @Param request body links.Link.Hash true "query params"
// @Success 200 {object} links.Link
// @Router /links/{key} [put]
func UpdateLinkByKey(ctx *gin.Context) {
	config.UpdateItemByParam(ctx, Cache, pkgName, Link{})
	return
}

// @Summary Delete link by its Key
// @Description delete link by its Key
// @Tags links
// @Produce json
// @Param  id  path  string  true  "link Key"
// @Success 200 {object} links.Link
// @Router /links/{key} [delete]
func DeleteLinkByKey(ctx *gin.Context) {
	config.DeleteItemByParam(ctx, Cache, pkgName)
	return
}

// @Summary Toggle active boolean for {hash}
// @Description toggle active boolean for {hash}
// @Tags links
// @Produce json
// @Param  id  path  string  true  "hash"
// @Success 200 {object} links.Link
// @Router /links/{key}/active [put]
func ActiveToggleByKey(ctx *gin.Context) {
	var hash string = ctx.Param("key")
	var link Link

	rawLink, ok := Cache.Get(hash)
	if !ok {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "link not found",
		})
		return
	}

	link, ok = rawLink.(Link)
	if !ok {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "cannot assert data type, database internal error",
		})
		return
	}

	// inverse the Active field value
	link.Active = !link.Active

	if saved := Cache.Set(hash, link); !saved {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "link couldn't be saved to database",
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "link active toggle pressed!",
		"link":    link,
	})
	return
}
