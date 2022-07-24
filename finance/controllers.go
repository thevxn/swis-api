package finance

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAccounts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"message":  "dumping finance accounts",
		"accounts": finance.Accounts,
	})
}

func GetAccountByOwner(c *gin.Context) {
	owner := c.Param("owner")

	for _, f := range finance.Accounts {
		if f.Owner == owner {
			c.IndentedJSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"account": f,
			})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "account not found",
	})
}

func PostDumpRestore(c *gin.Context) {
	var importFinance Finance

	if err := c.BindJSON(&importFinance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	finance = importFinance

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "finance imported, omitting output",
	})
}
