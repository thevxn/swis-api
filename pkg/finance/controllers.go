package finance

import (
	"net/http"
	"strconv"

	"go.savla.dev/swis/v5/pkg/core"

	"github.com/gin-gonic/gin"
)

var (
	Cache   *core.Cache
	pkgName string = "finance"
)

/*

  accounts

*/

func GetRootData(ctx *gin.Context) {}

// @Summary Get all finance accounts
// @Description get finance complete list
// @Tags finance
// @Produce json
// @Success 200 {object} finance.Account
// @Router /finance/accounts [get]
func GetAccounts(ctx *gin.Context) {
	core.PrintAllRootItems(ctx, Cache, pkgName)
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
	core.AddNewItemByParam(ctx, Cache, pkgName, Account{})
	return
}

// @Summary Get finance account by its Owner
// @Description get finance account by :owner param
// @Tags finance
// @Produce json
// @Success 200 {object} finance.Account
// @Router /finance/accounts/owner/{key} [get]
func GetAccountByOwnerKey(ctx *gin.Context) {
	//core.PrintItemByParam(ctx, Cache, pkgName, Account{})
	return
}

// @Summary Upload finance accounts dump backup -- restores all finance accounts
// @Description upload accounts JSON dump
// @Tags finance
// @Accept json
// @Produce json
// @Router /finance/restore [post]
func PostDumpRestore(ctx *gin.Context) {
	core.BatchRestoreItems(ctx, Cache, pkgName, Account{})
	return
}

// (PUT /finance/accounts/id})
// @Summary Update finance account by ID
// @Description update finance account by ID
// @Tags finance
// @Produce json
// @Param request body finance.Account.ID true "query params"
// @Success 200 {object} finance.Account
// @Router /finance/accounts/{key} [put]
func UpdateAccountByKey(ctx *gin.Context) {
	core.UpdateItemByParam(ctx, Cache, pkgName, Account{})
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
	core.DeleteItemByParam(ctx, Cache, pkgName)
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
	key := ctx.Param("owner")
	tax := Tax{}

	y := ctx.Param("year")
	year, err := strconv.Atoi(y)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"count":   0,
			"message": "invalid year on input",
			"tax":     tax,
			"key":     key,
		})
		return
	}

	//accounts := make(map[string]Account)
	accounts := []Account{}
	counter := 0

	// fetch users's account(s)
	rawAccounts, _ := Cache.GetAll()
	for key, rawAccount := range rawAccounts {
		acc, ok := rawAccount.(Account)
		if !ok {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "cannot assert type, backend data error",
				"tax":     tax,
				"key":     key,
			})
			return

			accounts = append(accounts, acc)
		}
	}

	// stop on zero accounts found
	if len(accounts) == 0 {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "no account for such owner",
			"tax":     tax,
			"key":     key,
		})
		return
	}

	// TODO un-hardcode this
	currency := "CZK"

	// iterate over accounts and items, do the taxes
	for _, acc := range accounts {
		// TODO we should consider that accounts could have different currencies set
		if acc.Currency != currency {
			continue
		}

		// calculate income and expenses
		for _, item := range acc.Items {
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

		tax.Summary += tax.IncomeTotal - tax.ExpenseTotal
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, sending year stats for the user's taxes",
		"tax":     tax,
		"key":     key,
		"year":    year,
	})
	return
}
