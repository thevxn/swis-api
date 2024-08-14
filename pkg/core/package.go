package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Package struct {
	Name        string
	Cache       []**Cache
	Routes      func(r *gin.RouterGroup)
	Generic     bool
	Subpackages []string
}

type FieldDetail struct {
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Readonly bool   `json:"readonly"`
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

func AddNewItem[T any](ctx *gin.Context, cache *Cache, pkgName string, model T) {
	// Read the body
	bodyBytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
			"message": "failed to read request body",
			"package": pkgName,
		})
		return
	}

	bodyCopy := bodyBytes

	meta := struct {
		ID string `json:"id"`
	}{}

	if err := json.Unmarshal(bodyCopy, &meta); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
			"message": "cannot determine the new ID",
			"package": pkgName,
		})
		return
	}

	key := meta.ID

	// Reset the body so it can be read again
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

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

func ParsePackageTypes(ctx *gin.Context, pkgName string, models ...interface{}) {
	var types = make(map[string]map[string]FieldDetail)

	for _, model := range models {
		typ := reflect.TypeOf(model)

		types[typ.Name()] = listFieldTypes(model)
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "parsing pkg's model types",
		"types":   types,
		"package": pkgName,
	})
	return
}

func ParsePackageType(ctx *gin.Context, pkgName string, model interface{}) {
	var types = make(map[string]FieldDetail)

	types = listFieldTypes(model)

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "parsing pkg's model field types",
		"types":   types,
		"package": pkgName,
	})
	return
}

func listFieldTypes(str interface{}) map[string]FieldDetail {
	var body = make(map[string]FieldDetail)

	val := reflect.ValueOf(str)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldType := field.Type

		jsonTag := field.Tag.Get("json")
		requiredTag := field.Tag.Get("required")
		roTag := field.Tag.Get("readonly")

		if fieldType.Kind() == reflect.Array || fieldType.Kind() == reflect.Slice {
			elemType := fieldType.Elem()
			if elemType.Kind() == reflect.Struct {
				body[jsonTag] = FieldDetail{
					Type:     "[]json",
					Required: requiredTag == "true",
					Readonly: roTag == "true",
				}
				continue
			}
		}

		if fieldType.Kind() == reflect.Struct || fieldType.Kind() == reflect.Map {
			body[jsonTag] = FieldDetail{
				Type:     "json",
				Required: requiredTag == "true",
				Readonly: roTag == "true",
			}
			continue
		}

		body[jsonTag] = FieldDetail{
			Type:     fmt.Sprintf("%s", fieldType),
			Required: requiredTag == "true",
			Readonly: roTag == "true",
		}
	}

	return body
}
