package swife

import (
	b64 "encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func findFrontendBySiteName(c *gin.Context) (f *Frontend) {
	for _, f := range swives {
		if f.SiteName == c.Param("sitename") {
			//c.IndentedJSON(http.StatusOK, a)
			return &f
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "frontend not found",
	})
	return nil
}

func GetFrontendBySiteName(c *gin.Context) {
	if frontend := findFrontendBySiteName(c); frontend != nil {

		// reencode HTML strings
		frontend.Title = b64.StdEncoding.
			EncodeToString([]byte(fmt.Sprintf("%s", frontend.Title)))
		frontend.Description = b64.StdEncoding.
			EncodeToString([]byte(fmt.Sprintf("%s", frontend.Description)))

		c.IndentedJSON(http.StatusOK, frontend)
	}
}

func PostDumpRestore(c *gin.Context) {
	var importFrontends []Frontend

	if err := c.BindJSON(&importFrontends); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	//depots = append(depots, importDepot)
	swives = importFrontends

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "swife frontends imported, omitting output",
	})
}
