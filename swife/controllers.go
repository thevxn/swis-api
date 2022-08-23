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

// @Summary Get frontend by Sitename
// @Description get frontend details by :sitename param
// @Tags swife
// @Produce  json
// @Success 200 {object} swife.Frontend
// @Router /swife/{sitename} [get]
func GetFrontendBySiteName(c *gin.Context) {
	if frontend := findFrontendBySiteName(c); frontend != nil {

		// reencode HTML strings --- only if raw HTML as input!
		//frontend.Title = b64.StdEncoding.
		//	EncodeToString([]byte(fmt.Sprintf("%s", frontend.Title)))
		//frontend.Description = b64.StdEncoding.
		//	EncodeToString([]byte(fmt.Sprintf("%s", frontend.Description)))
		c.IndentedJSON(http.StatusOK, frontend)
	}
}

// @Summary Uploadswife dump backup -- restores all frontends
// @Description upload frontend JSON dump and restore the data model
// @Tags swife
// @Accept json
// @Produce json
// @Router /swife/restore [post]
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
