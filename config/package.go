// generics and first class functions PoC
package config

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PrintAllRootItems(ctx *gin.Context, cache *Cache, dataName string) {
	items, count := cache.GetAll()

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"count":   count,
		"message": fmt.Sprintf("ok, listing all items of '%s' package", dataName),
		dataName:  items,
	})
	return
}

func PrintItemByParam[T any](ctx *gin.Context, cache *Cache, dataName string, model T) {
	var itemName string = ctx.Param("name")

	rawItem, ok := cache.Get(itemName)
	if !ok {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "item not found",
			"name":    itemName,
		})
		return
	}

	item, ok := rawItem.(T)
	if !ok {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "cannot assert data type, database internal error",
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, dumping item's contents",
		dataName:  item,
	})
	return
}

func AddNewItemByParam[T any](ctx *gin.Context, cache *Cache, dataName string, model T) {
	if err := ctx.BindJSON(&model); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot bind input JSON stream",
		})
		return
	}

	//key := model.Name | model.ID
	key := ctx.Param("name")

	// TODO: implement LoadOrStore() method
	if _, found := cache.Get(key); found {
		ctx.IndentedJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"key":     key,
			"message": "item already exists",
		})
		return
	}

	if saved := cache.Set(key, model); !saved {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "item couldn't be saved to database",
		})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"item":    model,
		"key":     key,
		"message": "new item added",
	})

	return
}

func UpdateItemByParam[T any](ctx *gin.Context, cache *Cache, dataName string, model T) {
	key := ctx.Param("name")

	if _, found := cache.Get(key); !found {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "item not found",
		})
		return
	}

	if err := ctx.BindJSON(&model); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot bind input JSON stream",
		})
		return
	}

	if saved := cache.Set(key, model); !saved {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "item couldn't be saved to database",
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"item":    model,
		"key":     key,
		"message": "item updated",
	})
	return
}

func DeleteItemByParam[T any](ctx *gin.Context, cache *Cache, dataName string, model T) {
	key := ctx.Param("name")

	if _, found := cache.Get(key); !found {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "item not found",
		})
		return
	}

	if deleted := cache.Delete(key); !deleted {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "item couldn't be deleted from database",
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"key":     key,
		"message": "item deleted by key",
	})
	return
}

func BatchRestoreItems[T any](ctx *gin.Context, cache *Cache, dataName string, model T) {
	var counter int = 0

	items := make(map[string]T)

	if err := ctx.BindJSON(&items); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot bind input JSON stream",
		})
		return
	}

	for key, item := range items {
		cache.Set(key, item)
		counter++
	}

	ctx.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"count":   counter,
		"message": "items restored successfully",
	})
	return
}
