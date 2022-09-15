package finance

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get all finance accounts
// @Description get finance complete list
// @Tags finance
// @Produce  json
// @Success 200 {object} finance.Finance
// @Router /finance [get]
func GetAccounts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"message":  "dumping finance accounts",
		"accounts": finance.Accounts,
	})
}

// @Summary Get finance account by its Owner
// @Description get finance account by :owner param
// @Tags finance
// @Produce  json
// @Success 200 {object} finance.Account
// @Router /finance/accounts/{owner} [get]
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

// @Summary Upload finance accounts dump backup -- restores all finance accounts
// @Description upload accounts JSON dump
// @Tags finance
// @Accept json
// @Produce json
// @Router /finance/restore [post]
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
