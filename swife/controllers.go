package swife

import (
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
		c.IndentedJSON(http.StatusOK, frontend)
	}
}
