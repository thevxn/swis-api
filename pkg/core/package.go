// generics and first class functions PoC
package core

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Package struct {
	Name    string
	Cache   []**Cache
	Routes  func(r *gin.RouterGroup)
	Generic bool
}

func PrintAllRootItems(ctx *gin.Context, cache *Cache, pkgName string) {
	items, count := cache.GetAll()

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"count":   count,
		"items":   items,
		"message": fmt.Sprintf("ok, listing all items of '%s' package", pkgName),
		"package": pkgName,
	})
	return
}

func PrintItemByParam[T any](ctx *gin.Context, cache *Cache, pkgName string, model T) {
	key := ctx.Param("key")

	rawItem, ok := cache.Get(key)
	if !ok {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"key":     key,
			"message": "item not found",
			"package": pkgName,
		})
		return
	}

	item, ok := rawItem.(T)
	if !ok {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"key":     key,
			"message": "cannot assert data type, database internal error",
			"package": pkgName,
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"item":    item,
		"key":     key,
		"message": "ok, dumping item's contents",
		"package": pkgName,
	})
	return
}

func AddNewItemByParam[T any](ctx *gin.Context, cache *Cache, pkgName string, model T) {
	//key := model.Name | model.ID
	key := ctx.Param("key")

	if err := ctx.BindJSON(&model); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
			"key":     key,
			"message": "cannot bind input JSON stream",
			"package": pkgName,
		})
		return
	}

	// TODO: implement LoadOrStore() method
	if _, found := cache.Get(key); found {
		ctx.IndentedJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"key":     key,
			"message": "item already exists",
			"package": pkgName,
		})
		return
	}

	if saved := cache.Set(key, model); !saved {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"key":     key,
			"message": "item couldn't be saved to database",
			"package": pkgName,
		})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"item":    model,
		"key":     key,
		"message": "new item added",
		"package": pkgName,
	})
	return
}

func UpdateItemByParam[T any](ctx *gin.Context, cache *Cache, pkgName string, model T) {
	key := ctx.Param("key")

	if _, found := cache.Get(key); !found {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"key":     key,
			"message": "item not found",
			"package": pkgName,
		})
		return
	}

	if err := ctx.BindJSON(&model); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
			"key":     key,
			"message": "cannot bind input JSON stream",
			"package": pkgName,
		})
		return
	}

	if saved := cache.Set(key, model); !saved {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"key":     key,
			"message": "item couldn't be saved to database",
			"package": pkgName,
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"item":    model,
		"key":     key,
		"message": "item updated",
		"packege": pkgName,
	})
	return
}

func DeleteItemByParam(ctx *gin.Context, cache *Cache, pkgName string) {
	key := ctx.Param("key")

	if _, found := cache.Get(key); !found {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"key":     key,
			"message": "item not found",
			"package": pkgName,
		})
		return
	}

	if deleted := cache.Delete(key); !deleted {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"key":     key,
			"message": "item couldn't be deleted from database",
			"package": pkgName,
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"key":     key,
		"message": "item deleted by key",
		"package": pkgName,
	})
	return
}

func BatchRestoreItems[T any](ctx *gin.Context, cache *Cache, pkgName string) {
	var counter int = 0

	items := struct {
		Items map[string]T `json:"items"`
	}{}

	if err := ctx.BindJSON(&items); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
			"message": "cannot bind input JSON stream",
			"package": pkgName,
		})
		return
	}

	for key, item := range items.Items {
		if key == "" {
			continue
		}
		cache.Set(key, item)
		counter++
	}

	ctx.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"count":   counter,
		"message": "items restored successfully",
		"package": pkgName,
	})
	return
}
