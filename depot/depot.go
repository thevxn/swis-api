package depot

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

type Depots struct {
	Depots  []Depot `json:"depots"`
}

type Depot struct {
	Owner   	string	`json:"owner_name"`
	DepotItems 	[]Item 	`json:"depot_items"`
}

type Item struct {
	ID       	int 	`json:"id"`
	Description 	string 	`json:"desc"`
	Misc     	string 	`json:"misc"`
	Location	string 	`json:"depot"`
}

// flush depots at start
var depots = Depots{}


func GetDepots(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message": "dumping depots",
		"depots": depots,
	})
}

func GetDepotByOwner(c *gin.Context) {
	owner := c.Param("owner")

	for _, d := range depots.Depots {
		if d.Owner == owner {
			// https://pkg.go.dev/sort#Slice
			sort.Slice(d.DepotItems, func(i, j int) bool {
    				return (d.DepotItems[i].Description < d.DepotItems[j].Description)
			})

			c.IndentedJSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"message": "dumping user's depot, alphabetically sorted",
				"depot": d,
			})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code": http.StatusNotFound,
		"message": "depot not found",
	})
}


func PostDumpRestore(c *gin.Context) {
	var importDepots Depots

	if err := c.BindJSON(&importDepots); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	//depots = append(depots, importDepot)
	depots = importDepots

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated,
		"message": "depots imported, omitting output",
		//"depots": importDepots.Depots,
	})
}

