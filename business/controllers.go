package business

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func findBusinessByICO(c *gin.Context) (b *Business) {
	// loop over businesses
	for _, b := range businessArray.BusinessArray {
		if b.ICO == c.Param("ico_id") {
			//c.IndentedJSON(http.StatusOK, b)
			return &b
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "business not found by given ICO ID",
	})
	return nil
}

func GetBusinessArray(c *gin.Context) {
	// serialize struct to JSON
	// net/http response code
	//c.IndentedJSON(http.StatusOK, businessArray)
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"message":  "dumping business array",
		"business": businessArray.BusinessArray,
	})
}

// GetBusinessByICO returns business' properties, given sent ICO exists in database.
func GetBusinessByICO(c *gin.Context) {
	if biz := findBusinessByICO(c); biz != nil {
		// business found
		c.IndentedJSON(http.StatusOK, biz)
	}
}

// PostBusiness
func PostBusiness(c *gin.Context) {
	var newBusiness Business

	// bind received JSON to newUser
	if err := c.BindJSON(&newBusiness); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	// add new business
	businessArray.BusinessArray = append(businessArray.BusinessArray, newBusiness)
	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, newBusiness)
}

// PostDumpRestore
func PostDumpRestore(c *gin.Context) {
	var importBusinessArray BusinessArray

	if err := c.BindJSON(&importBusinessArray); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	businessArray = importBusinessArray

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "business array imported, omitting output",
	})
}