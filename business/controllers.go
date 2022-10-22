package business

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var b sync.Map

// @Summary Get all business entities
// @Description get business entities list
// @Tags business
// @Produce  json
// @Success 200 {object} business.Entities
// @Router /business [get]
func GetBusinessEntities(c *gin.Context) {
	var entities = make(map[string]Business)

	b.Range(func(rawKey, rawVal interface{}) bool {
		k, ok := rawKey.(string)
		v, ok := rawVal.(Business)

		if !ok {
			return false
		}

		entities[k] = v
		return true
	})

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"message":  "ok, listing business entities",
		"business": entities,
	})
	return
}

// @Summary Get business entity by its ID
// @Description get business by ID param
// @Tags business
// @Produce  json
// @Success 200 {object} business.Business
// @Router /business/{id} [get]
func GetBusinessByID(c *gin.Context) {
	var id string = c.Param("id")
	var entity Business

	rawEntity, ok := b.Load(id)
	entity, ok = rawEntity.(Business)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "business entity not found",
			"code":    http.StatusNotFound,
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message":  "ok, dumping requested business entity's details",
		"code":     http.StatusOK,
		"business": entity,
	})
	return
}

// @Summary Add new business entity
// @Description add new business entity
// @Tags business
// @Produce json
// @Param request body business.Business true "query params"
// @Success 200 {object} business.Business
// @Router /business [post]
func PostBusiness(c *gin.Context) {
	var newBusiness = &Business{}

	if err := c.BindJSON(newBusiness); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	var id string = newBusiness.ID

	if _, found := b.Load(id); found {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "business entity's ID already used!",
			"id":      id,
		})
		return
	}

	b.Store(newBusiness.ID, newBusiness)

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":     http.StatusCreated,
		"message":  "business entity added",
		"id":       id,
		"business": newBusiness,
	})
	return
}

// (PUT /business/{id})
// @Summary Update business entity by its ID
// @Description update business entity by its ID
// @Tags business
// @Produce json
// @Param request body business.Business.ID true "query params"
// @Success 200 {object} business.Business
// @Router /business/{id} [put]
func UpdateBusinessByID(c *gin.Context) {
	var id string = c.Param("id")
	var updatedEntity Business

	if _, ok := b.Load(id); !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "business entity not found",
			"code":    http.StatusNotFound,
		})
		return
	}

	if err := c.BindJSON(&updatedEntity); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	b.Store(id, updatedEntity)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"message":  "business entity updated",
		"business": updatedEntity,
	})
	return
}

// @Summary Delete business by its ID
// @Description delete business by its ID
// @Tags business
// @Produce json
// @Param  id  path  string  true  "business ID"
// @Success 200 {object} business.Business.ID
// @Router /business/{id} [delete]
func DeleteBusinessByID(c *gin.Context) {
	var id string = c.Param("id")

	b.Delete(id)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "business deleted by ID",
		"id":      id,
	})
	return
}

// @Summary Upload business dump backup -- restores all business entities
// @Description upload business JSON dump
// @Tags business
// @Accept json
// @Produce json
// @Router /business/restore [post]
// PostDumpRestore
func PostDumpRestore(c *gin.Context) {
	var importEntities = &Entities{}

	if err := c.BindJSON(importEntities); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	for _, business := range importEntities.Entities {
		b.Store(business.ID, business)
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"message": "business entities imported, omitting output",
		"code":    http.StatusCreated,
	})
	return
}
