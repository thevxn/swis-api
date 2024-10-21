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

// Package struct describes the structure of a generic package to be loaded into the engine at start.
type Package struct {
	// Name hold the name of a package.
	Name string

	// Cache is an array of pointers to caches to be initialized.
	Cache []**Cache

	// CacheName is an array of names for such caches being initialized.
	CacheNames []string

	// Routes is a function which holds the package's routes with their methods specified too.
	Routes func(r *gin.RouterGroup)

	// Generic is a boolean indicating whether is the root package a generic CRUD package (does not contain any subpackages).
	Generic bool

	// Subpackages is an array of subpackage names to register as generic ones.
	Subpackages []string

	// SubpackageModels is a map to match the root model for such subpackage.
	SubpackageModels map[string]any
}

// FieldDetail is a struct to describe any loaded model's field for the type enum export.
type FieldDetail struct {
	// Type holds the name of suck type.
	Type string `json:"type"`

	// Required indicates whether is the field required or not.
	Required bool `json:"required"`

	// ReadOnly indicates whether is the field readonly (write-locked) or not.
	Readonly bool `json:"readonly"`
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

func BatchRestoreItems[T any](ctx *gin.Context, pkg Package) {
	var counter int = 0

	if len(pkg.Subpackages) == 0 {
		items := struct {
			Items map[string]T `json:"items"`
		}{}

		if err := ctx.BindJSON(&items); err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"error":   err.Error(),
				"message": "cannot bind input JSON stream",
				"package": pkg.Name,
			})
			return
		}

		cache := *pkg.Cache[0]

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
			"package": pkg.Name,
		})
		return
	}

	if len(pkg.Subpackages) != len(pkg.SubpackageModels) {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot restore data: wrong package configuration",
			"package": pkg.Name,
		})
		return
	}

	for idx, subpkg := range pkg.Subpackages {
		// Otherwise, we have to extract each subpackage's data manually.
		items := subpkgItems{}

		if err := ctx.BindJSON(&items); err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"error":   err.Error(),
				"message": "cannot bind input JSON stream (subpackages)",
				"package": pkg.Name,
			})
			return
		}

		// Assert subpackage's model type
		var ok bool
		counter, ok = restoreSubpackageData(items, pkg.SubpackageModels, subpkg, *pkg.Cache[idx])
		if !ok {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "cannot restore data: cannot assert subpackage's model type",
				"package": pkg.Name,
			})
			return
		}
	}

	ctx.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"count":   counter,
		"message": "items restored successfully (subpackages)",
		"package": pkg.Name,
	})
	return
}

type subpkgItems struct {
	Items map[string]interface{} `json:"items"`
}

func restore[T any](input interface{}, model T) (map[string]T, bool) {
	m, ok := input.(map[string]T)
	if !ok {
		return nil, false
	}

	return m, true
}

func restoreSubpackageData[K string, V any](items subpkgItems, m map[string]V, subpkg string, cache *Cache) (int, bool) {
	var counter int

	for k, v := range m {
		if k != subpkg {
			continue
		}

		subData, ok := restore(items.Items[subpkg], v)
		if !ok {
			return 0, false
		}

		for key, item := range subData {
			if key == "" {
				continue
			}

			ok = cache.Set(key, item)
			if !ok {
				return 0, false
			}
			counter++
		}
	}

	return counter, true
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
