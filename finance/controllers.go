package finance

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var f sync.Map

// @Summary Get all finance accounts
// @Description get finance complete list
// @Tags finance
// @Produce json
// @Success 200 {object} finance.Account
// @Router /finance [get]
func GetAccounts(c *gin.Context) {
	var accounts = make(map[string]Account)

	f.Range(func(rawKey, rawVal interface{}) bool {
		// very insecure assert
		k, ok := rawKey.(string)
		v, ok := rawVal.(Account)

		if !ok {
			return false
		}

		accounts[k] = v
		return true
	})

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"message":  "dumping finance accounts",
		"accounts": accounts,
	})
}

// @Summary Get finance account by its Owner
// @Description get finance account by :owner param
// @Tags finance
// @Produce json
// @Success 200 {object} finance.Account
// @Router /finance/accounts/{owner} [get]
func GetAccountByOwner(c *gin.Context) {
	var accounts = make(map[string]Account)
	var owner string = c.Param("owner")

	f.Range(func(rawKey, rawVal interface{}) bool {
		// very insecure assert
		k, ok := rawKey.(string)
		v, ok := rawVal.(Account)

		if !ok {
			return false
		}

		if v.Owner == owner {
			accounts[k] = v
		}
		return true
	})

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"message":  "dumping accouts by owner",
		"accounts": accounts,
	})
	return
}

// @Summary Upload finance accounts dump backup -- restores all finance accounts
// @Description upload accounts JSON dump
// @Tags finance
// @Accept json
// @Produce json
// @Router /finance/restore [post]
func PostDumpRestore(c *gin.Context) {
	var importAccounts = &Accounts{}

	if err := c.BindJSON(importAccounts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	for _, acc := range importAccounts.Accounts {
		f.Store(acc.ID, acc)
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "finance accounts imported, omitting output",
	})
}
