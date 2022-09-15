package six

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// (GET /six)
// @Summary Get the six struct
// @Description get the six struct
// @Tags six
// @Produce  json
// @Success 200 {object} six.SixStruct
// @Router /six [get]
func GetSixStruct(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, dumping the six struct",
		"six":     sixStruct,
	})
}
