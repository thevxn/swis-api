package finance

import (
	//"net/http"

	"go.savla.dev/swis/v5/config"

	"github.com/gin-gonic/gin"
)

var (
	Cache   *config.Cache
	pkgName string = "finance"
)

// @Summary Get all finance accounts
// @Description get finance complete list
// @Tags finance
// @Produce json
// @Success 200 {object} finance.Account
// @Router /finance/account [get]
func GetAccounts(ctx *gin.Context) {
	config.PrintAllRootItems(ctx, Cache, pkgName)
	return
}

// @Summary Add new finance account
// @Description add new finance account
// @Tags finance
// @Produce json
// @Param request body finance.Account true "query params"
// @Success 200 {object} finance.Account
// @Router /finance/account/{key} [post]
func PostNewAccountByKey(ctx *gin.Context) {
	config.AddNewItemByParam(ctx, Cache, pkgName, Account{})
	return
}

// @Summary Get finance account by its Owner
// @Description get finance account by :owner param
// @Tags finance
// @Produce json
// @Success 200 {object} finance.Account
// @Router /finance/account/owner/{key} [get]
func GetAccountByOwnerKey(ctx *gin.Context) {
	//config.PrintItemByParam(ctx, Cache, pkgName, Account{})
	return
}

// @Summary Upload finance accounts dump backup -- restores all finance accounts
// @Description upload accounts JSON dump
// @Tags finance
// @Accept json
// @Produce json
// @Router /finance/restore [post]
func PostDumpRestore(ctx *gin.Context) {
	config.BatchRestoreItems(ctx, Cache, pkgName, Account{})
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
	config.UpdateItemByParam(ctx, Cache, pkgName, Account{})
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
	config.DeleteItemByParam(ctx, Cache, pkgName)
	return
}
