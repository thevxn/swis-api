package news

import (
	"net/http"
	"net/url"
	"sort"
	"time"

	"go.vxn.dev/swis/v5/pkg/core"

	"github.com/gin-gonic/gin"
)

var (
	Cache *core.Cache

	caches = []**core.Cache{
		&Cache,
	}
	pkgName string = "news"
)

var Package *core.Package = &core.Package{
	Name:   pkgName,
	Cache:  caches,
	Routes: Routes,
	Subpackages: []string{
		"sources",
	},
}

var restorePackage = &core.RestorePackage{
	Name:             pkgName,
	Cache:            caches,
	CacheNames:       []string{"Cache"},
	Subpackages:      []string{},
	SubpackageModels: map[string]any{},
}

// GetSources
// @Summary Get news source list
// @Description get all news sources
// @Tags news
// @Produce  json
// @Success 200 {object} news.NewsSources.Sources
// @Router /news/sources [get]
func GetSources(ctx *gin.Context) {
	core.PrintAllRootItems(ctx, Cache, pkgName)
	return
}

// GetSourcesByUserKey
// @Summary Get news source list by User key
// @Description get news sources by their user :key param
// @Tags news
// @Produce  json
// @Success 200 {object} news.UserSource
// @Router /news/sources/{key} [get]
func GetSourcesByUserKey(ctx *gin.Context) {
	core.PrintItemByParam[UserSource](ctx, Cache, pkgName, UserSource{})
	return
}

// @Summary Add new user sources
// @Description add new news sources
// @Tags news
// @Produce  json
// @Success 200 {object} news.UserSource
// @Router /news/sources [post]
func PostNewSources(ctx *gin.Context) {
	core.AddNewItem[UserSource](ctx, Cache, pkgName, UserSource{})
	return
}

// @Summary Update news sources by user key
// @Description update news sources by user key
// @Tags news
// @Produce json
// @Param request body news.Source true "query params"
// @Success 200 {object} news.UserSource
// @Router /news/sources/{key} [put]
func UpdateSourcesByUserKey(ctx *gin.Context) {
	core.UpdateItemByParam[UserSource](ctx, Cache, pkgName, UserSource{})
	return
}

// @Summary Delete user sources by user key
// @Description delete user sources by user key
// @Tags news
// @Produce json
// @Success 200 {string} string "ok"
// @Router /news/sources/{key} [delete]
func DeleteSourcesByUserKey(ctx *gin.Context) {
	core.DeleteItemByParam(ctx, Cache, pkgName)
	return
}

// PostDumpRestore
// @Summary Upload news sources dump backup -- restores all sources
// @Description update news sources JSON dump
// @Tags news
// @Accept json
// @Produce json
// @Router /news/sources/restore [post]
func PostDumpRestore(ctx *gin.Context) {
	core.BatchRestoreItems[UserSource](ctx, restorePackage)
	return
}

// @Summary List package model's field types
// @Description list package model's field types
// @Tags news
// @Accept json
// @Produce json
// @Router /news/sources/types [get]
func ListTypesSources(ctx *gin.Context) {
	core.ParsePackageType(ctx, pkgName, UserSource{})
	return
}

// GetNewsByUser returns all possible news from all sources loaded in memory
// @Summary Get news by user key
// @Description fetch and parse news for :key param
// @Tags news
// @Produce  json
// @Success 200 {object} news.Item
// @Router /news/{key} [get]
func GetNewsByUserKey(ctx *gin.Context) {
	user := ctx.Param("key")

	rawUserSources, ok := Cache.Get(user)
	if !ok {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"key":     user,
			"message": "no news sources found for such user key",
			"package": pkgName,
		})
		return
	}

	userSources, ok := rawUserSources.(UserSource)
	if !ok {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"key":     user,
			"message": "cannot assert data type, internal database error",
			"package": pkgName,
		})
		return
	}

	//var R = []Rss{}
	var items = []Item{}

	for _, source := range userSources.Sources {
		contents := fetchRSSContents(source)
		if contents == nil {
			continue
		}

		for _, item := range *contents {
			// time layouts (date template constants) --> https://go.dev/src/time/format.go
			item.ParseDate, _ = time.Parse(time.RFC1123Z, item.PubDate)

			// convert news link to server/hostname
			u, _ := url.Parse(item.Link)
			item.Server = string(u.Hostname())

			items = append(items, item)
		}
	}

	// sort items by date DESC
	// https://stackoverflow.com/a/47028486
	sort.Slice(items, func(i, j int) bool {
		return items[i].ParseDate.After(items[j].ParseDate)
	})

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"items":   items,
		"message": "ok, listing news (newest to oldest)",
		"package": pkgName,
	})
	return
}
