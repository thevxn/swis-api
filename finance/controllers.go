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

// @Summary Add new finance account
// @Description add new finance account
// @Tags finance
// @Produce json
// @Param request body finance.Account true "query params"
// @Success 200 {object} finance.Account
// @Router /finance [post]
func PostNewAccount(c *gin.Context) {
	var newAccount *Account = &Account{}

	if err := c.BindJSON(newAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	var id string = newAccount.ID

	if _, found := f.Load(id); found {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "account ID already used!",
			"id":      id,
		})
		return
	}

	f.Store(newAccount.ID, newAccount)

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "link added",
		"id":      id,
		"account": newAccount,
	})
	return
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

// (PUT /finance/accounts/id})
// @Summary Update finance account by ID
// @Description update finance account by ID
// @Tags finance
// @Produce json
// @Param request body finance.Account.ID true "query params"
// @Success 200 {object} finance.Account
// @Router /finance/accounts/{id} [put]
func UpdateAccountByID(c *gin.Context) {
	var account *Account = &Account{}
	var id string = c.Param("id")

	if _, ok := f.Load(id); !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "account not found",
			"code":    http.StatusNotFound,
		})
		return
	}

	if err := c.BindJSON(account); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	f.Store(id, account)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "socket updated",
		"account": account,
	})
	return
}

// @Summary Delete finance account by ID
// @Description delete finance account by ID
// @Tags finance
// @Produce json
// @Param  id  path  string  true  "account ID"
// @Success 200 {object} finance.Account
// @Router /finance/accounts/{id} [delete]
func DeleteAccountByID(c *gin.Context) {
	var id string = c.Param("id")

	f.Delete(id)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "account deleted by Hash",
		"id":      id,
	})
	return
}
