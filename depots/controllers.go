package depots

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

// @Summary Get all depots and their users/owners
// @Description get depot complete list
// @Tags depot
// @Produce json
// @Success 200 {object} depots.Depots
// @Router /depots [get]
// GetSocketList GET method
func GetDepots(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, dumping depots",
		"depots":  depots.Depots,
	})
}

// @Summary Get depot list by Owner
// @Description get depot list by :owner param
// @Tags depot
// @Produce json
// @Success 200 {object} depots.Depot
// @Router /depots/{owner} [get]
// GetSocketList GET method
func GetDepotByOwner(c *gin.Context) {
	owner := c.Param("owner")

	for _, depot := range depots.Depots {
		if depot.Owner == owner {
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
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "depot not found",
	})
}

// @Summary Upload depot dump backup -- restores all depots
// @Description upload depot JSON dump
// @Tags depot
// @Accept json
// @Produce json
// @Router /depots/restore [post]
func PostDumpRestore(c *gin.Context) {
	var importDepots Depots

	if err := c.BindJSON(&importDepots); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	//depots = append(depots, importDepot)
	depots = importDepots

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "depots imported, omitting output",
	})
}
