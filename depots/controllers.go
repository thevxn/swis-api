package depots

import (
	"net/http"
	"sort"
	"sync"

	"github.com/gin-gonic/gin"
)

var d sync.Map

// @Summary Get all depots and their users/owners
// @Description get depot complete list
// @Tags depots
// @Produce json
// @Success 200 {object} depots.Depots
// @Router /depots [get]
// GetSocketList GET method
func GetDepots(c *gin.Context) {
	var depots = make(map[string]Depot)

	d.Range(func(rawKey, rawVal interface{}) bool {
		// very insecure assert
		k, ok := rawKey.(string)
		v, ok := rawVal.(Depot)

		if !ok {
			return false
		}

		depots[k] = v
		return true
	})
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, dumping depots",
		"depots":  depots,
	})
}

// @Summary Add new depot
// @Description add new depot
// @Tags depots
// @Produce json
// @Param request body depots.Depot true "query params"
// @Success 200 {object} depots.Depot
// @Router /depots [post]
func PostNewDepot(c *gin.Context) {
	var newDepot = &Depot{}

	if err := c.BindJSON(newDepot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	var owner string = newDepot.Owner

	_, found := d.Load(owner)
	if found {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "depot by owner name already created!",
			"owner":   owner,
		})
		return
	}

	d.Store(newDepot.Owner, newDepot)

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "depot added",
		"owner":   owner,
	})
}

// @Summary Get depot list by Owner
// @Description get depot list by :owner param
// @Tags depots
// @Produce json
// @Success 200 {object} depots.Depot
// @Router /depots/{owner} [get]
func GetDepotByOwner(c *gin.Context) {
	var depot Depot
	var owner string = c.Param("owner")

	rawDepot, ok := d.Load(owner)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "depot not found",
		})
		return
	}

	// assert type Depot
	depot, ok = rawDepot.(Depot)
	if !ok {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "loaded value not of type Depot",
		})
		return
	}

	// order items ASC alphabetically
	// https://pkg.go.dev/sort#Slice
	sort.Slice(depot.DepotItems, func(i, j int) bool {
		return (depot.DepotItems[i].Description < depot.DepotItems[j].Description)
	})

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, dumping user's depot, alphabetically sorted",
		"depot":   depot,
	})
	return
}

// (DELETE /depot/{owner})
// @Summary Delete depot by its owner
// @Description delete depot by its Owner
// @Tags depots
// @Produce json
// @Param  id  path  string  true  "depot Owner"
// @Success 200 {object} depots.Depot
// @Router /depots/{owner} [delete]
func DeleteDepotByOwner(c *gin.Context) {
	var owner string = c.Param("owner")

	if owner == "" {
		return
	}

	d.Delete(owner)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "depot deleted by Owner",
		"owner":   owner,
	})
	return
}

// @Summary Upload depot dump backup -- restores all depots
// @Description upload depot JSON dump
// @Tags depots
// @Accept json
// @Produce json
// @Router /depots/restore [post]
func PostDumpRestore(c *gin.Context) {
	var importDepots = &Depots{}

	if err := c.BindJSON(importDepots); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	for _, depot := range importDepots.Depots {
		d.Store(depot.Owner, depot)
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "depots imported, omitting output",
	})
}
