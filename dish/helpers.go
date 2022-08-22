package dish

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// contains checks if a string is present in a slice
// https://freshman.tech/snippets/go/check-if-slice-contains-element/
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func findSocketByID(c *gin.Context) (index *int, s *Socket) {
	for i, s := range socketArray {
		if s.ID == c.Param("id") {
			//c.IndentedJSON(http.StatusOK, a)
			return &i, &s
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "socket not found",
	})
	return nil, nil
}
