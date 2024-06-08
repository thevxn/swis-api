package finance

import (
	"net/http"
	"strconv"

	"go.savla.dev/swis/v5/pkg/core"

	"github.com/gin-gonic/gin"
)

var (
	CacheAccounts *core.Cache
	CacheItems    *core.Cache
	pkgName       string = "finance"
)

var Package *core.Package = &core.Package{
	Name: pkgName,
	Cache: []**core.Cache{
		&CacheAccounts,
		&CacheItems,
	},
	Routes: Routes,
}

/*

  accounts

*/

// @Summary Get all finance accounts
// @Description get finance complete list
// @Tags finance
// @Produce json
// @Success 200 {object} finance.Account
// @Router /finance/accounts [get]
func GetAccounts(ctx *gin.Context) {
	core.PrintAllRootItems(ctx, CacheAccounts, pkgName)
	return
}

// @Summary Add new finance account
// @Description add new finance account
// @Tags finance
// @Produce json
// @Param request body finance.Account true "query params"
// @Success 200 {object} finance.Account
// @Router /finance/accounts/{key} [post]
func PostNewAccountByKey(ctx *gin.Context) {
	core.AddNewItemByParam[Account](ctx, CacheAccounts, pkgName)
	return
}

// @Summary Get finance account by its key
// @Description get finance account by its key
// @Tags finance
// @Produce json
// @Success 200 {object} finance.Account
// @Router /finance/accounts/{key} [get]
func GetAccountByKey(ctx *gin.Context) {
	core.PrintItemByParam[Account](ctx, CacheAccounts, pkgName)
	return
}

// @Summary Update finance account by ID
// @Description update finance account by ID
// @Tags finance
// @Produce json
// @Param request body finance.Account.ID true "query params"
// @Success 200 {object} finance.Account
// @Router /finance/accounts/{key} [put]
func UpdateAccountByKey(ctx *gin.Context) {
	core.UpdateItemByParam[Account](ctx, CacheAccounts, pkgName)
	return
}

// @Summary Delete finance account by ID
// @Description delete finance account by ID
// @Tags finance
// @Produce json
// @Param  id  path  string  true  "account ID"
// @Success 200 {object} finance.Account
// @Router /finance/accounts/{key} [delete]
func DeleteAccountByKey(ctx *gin.Context) {
	core.DeleteItemByParam(ctx, CacheAccounts, pkgName)
	return
}

// @Summary Get finance accounts by Owner key
// @Description get finance accounts by Owner key
// @Tags finance
// @Produce json
// @Success 200 {object} finance.Account
// @Router /finance/accounts/owner/:key [get]
func GetAccountByOwnerKey(ctx *gin.Context) {
	owner := ctx.Param("key")
	exportedAccounts := make(map[string]Account)
	counter := 0

	rawAccounts, _ := CacheAccounts.GetAll()

	for _, rawAccount := range rawAccounts {
		acc, ok := rawAccount.(Account)
		if !ok {
			continue
		}

		if acc.Owner == owner {
			exportedAccounts[acc.ID] = acc
			counter++
		}
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"count":   counter,
		"items":   exportedAccounts,
		"message": "ok, dumping accounts by owner",
		"owner":   owner,
	})
	return
}

/*


  items

*/

// @Summary Get all account items
// @Description get account list of items
// @Tags finance
// @Produce json
// @Success 200 {object} finance.Item
// @Router /finance/items [get]
func GetItems(ctx *gin.Context) {
	core.PrintAllRootItems(ctx, CacheItems, pkgName)
	return
}

// @Summary Add new account item
// @Description add new account item
// @Tags finance
// @Produce json
// @Param request body finance.Item true "query params"
// @Success 200 {object} finance.Item
// @Router /finance/items/{key} [post]
func PostNewItemByKey(ctx *gin.Context) {
	core.AddNewItemByParam[Item](ctx, CacheItems, pkgName)
	return
}

// @Summary Get account item by key
// @Description get account item by its key
// @Tags finance
// @Produce json
// @Success 200 {object} finance.Item
// @Router /finance/items/{key} [get]
func GetItemByKey(ctx *gin.Context) {
	core.PrintItemByParam[Item](ctx, CacheItems, pkgName)
	return
}

// @Summary Update account item by its key
// @Description update account item ba its key
// @Tags finance
// @Produce json
// @Param request body finance.Item.ID true "query params"
// @Success 200 {object} finance.Item
// @Router /finance/items/{key} [put]
func UpdateItemByKey(ctx *gin.Context) {
	core.UpdateItemByParam[Item](ctx, CacheItems, pkgName)
	return
}

// @Summary Delete account item by its ID
// @Description delete account item by its ID
// @Tags finance
// @Produce json
// @Param  id  path  string  true  "item ID"
// @Success 200 {object} finance.Item
// @Router /finance/items/{key} [delete]
func DeleteItemByKey(ctx *gin.Context) {
	core.DeleteItemByParam(ctx, CacheItems, pkgName)
	return
}

// @Summary Get account items by account ID
// @Description get account items by account ID
// @Tags finance
// @Produce json
// @Success 200 {object} finance.Item
// @Router /finance/items/account/:key [get]
func GetItemsByAccountID(ctx *gin.Context) {
	acc := ctx.Param("key")
	exportedItems := make(map[string]Item)
	counter := 0

	rawItems, _ := CacheItems.GetAll()

	for _, rawItem := range rawItems {
		item, ok := rawItem.(Item)
		if !ok {
			continue
		}

		if item.AccountID == acc {
			exportedItems[item.ID] = item
			counter++
		}
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"count":   counter,
		"items":   exportedItems,
		"message": "ok, dumping items by accountID",
		"account": acc,
	})
	return
}

/*

  taxes

*/

// @Summary Do taxes by the account owner Owner key
// @Description do taxes by the account Owner key
// @Tags finance
// @Produce json
// @Param  id  path  string  true  "owner key"
// @Success 200 {object} finance.Tax
// @Router /finance/taxes/{owner}/{year} [get]
func DoTaxesByOwner(ctx *gin.Context) {
	owner := ctx.Param("owner")
	tax := Tax{}
	counter := 0

	y := ctx.Param("year")
	year, err := strconv.Atoi(y)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "invalid year on input",
			"tax":     tax,
			"key":     owner,
		})
		return
	}

	accounts := make(map[string]Account)
	items := make(map[string][]Item)
	keys := []string{}

	// fetch users's account(s)
	rawAccounts, _ := CacheAccounts.GetAll()
	for key, rawAccount := range rawAccounts {
		acc, ok := rawAccount.(Account)
		if !ok {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "cannot assert type, backend data error",
				"tax":     tax,
				"account": key,
				"key":     owner,
			})
			return
		}

		if acc.Owner == owner {
			accounts[acc.ID] = acc
			keys = append(keys, key)
		}
	}

	// fetch all items and filter out only those for gotten accounts
	rawItems, _ := CacheItems.GetAll()
	for key, rawItem := range rawItems {
		item, ok := rawItem.(Item)
		if !ok {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "cannot assert type, backend data error",
				"tax":     tax,
				"item":    key,
				"key":     owner,
			})
			return
		}

		if contains(keys, item.AccountID) {
			items[item.AccountID] = append(items[item.AccountID], item)
		}
	}

	// TODO un-hardcode this
	currency := "CZK"

	// stop on zero accounts found
	if len(accounts) == 0 {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"code":     http.StatusOK,
			"message":  "no account for such owner",
			"tax":      tax,
			"key":      owner,
			"accounts": keys,
		})
		return
	}

	// iterate over account keys --- accounts IDs
	for _, id := range keys {
		acc := accounts[id]

		// TODO we should consider that accounts could have different currencies set
		if acc.Currency != currency {
			continue
		}

		// calculate income and expenses
		for _, item := range items[id] {
			// TODO we should consider that accounts could have different currencies set
			if item.Currency != currency {
				continue
			}

			// filter out items dated differently (different year)
			if item.PaymentDate.Year() != year {
				continue
			}

			// do the initial calculation
			switch item.Type {
			case "income":
				tax.IncomeTotal += item.Amount
				counter++
				break

			case "expense":
				tax.ExpenseTotal += item.Amount
				counter++
				break

			default:
				continue
			}
		}
	}

	// define the tax constants
	const incTaxRate = 0.15
	const expRate60 = 0.60

	// calculate the base for the income tax
	tax.Base = tax.IncomeTotal - tax.ExpenseTotal

	// calculate two estimations of the income tax
	tax.IncomeTaxEstAbs = tax.Base * incTaxRate
	tax.IncomeTaxEst60 = tax.IncomeTotal * (1 - expRate60) * incTaxRate

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"message":  "ok, sending year stats for the user's taxes estimation",
		"tax":      tax,
		"key":      owner,
		"year":     year,
		"currency": currency,
	})
	return
}

/*

  restoration

*/

// @Summary Get whole finance package content
// @Description get whole finance package content
// @Tags finance
// @Produce  json
// @Router /finance [get]
func GetRootData(ctx *gin.Context) {
	accounts, _ := CacheAccounts.GetAll()
	items, _ := CacheItems.GetAll()

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"message":  "ok, dumping dish root",
		"accounts": accounts,
		"items":    items,
	})
}

// @Summary Upload account iteme dump backup -- restores all finance accounts
// @Description upload accounts JSON dump
// @Tags finance
// @Accept json
// @Produce json
// @Router /finance/restore [post]
func PostDumpRestore(ctx *gin.Context) {
	var counter []int = []int{0, 0}

	var importFinance = struct {
		Accounts map[string]Account `json:"accounts"`
		Items    map[string]Item    `json:"items"`
	}{}

	if err := ctx.BindJSON(&importFinance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
			"message": "cannot parse input JSON stream",
		})
		return
	}

	for key, item := range importFinance.Accounts {
		if key == "" {
			continue
		}

		CacheAccounts.Set(key, item)
		counter[0]++
	}

	for key, item := range importFinance.Items {
		if key == "" {
			continue
		}

		CacheItems.Set(key, item)
		counter[1]++
	}

	// HTTP 201 Created
	ctx.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"counter": counter,
		"message": "finance dump imported successfully",
	})
}
