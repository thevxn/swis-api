package finance

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Finance struct {
	Accounts  []Account `json:"accounts"`
}

type Account struct {
	Name	   	string	`json:"account_name"`
	AccountNumber	string  `json:"account_number"`
	// bank codes like "0100" would be invalid as type int
	BankCode	string	`json:"bank_code"`
	SWIFT		string	`json:"account_swift"`
	IBAN		string	`json:"account_iban"`
	Owner   	string	`json:"account_owner"`
	Items 		[]Item 	`json:"account_items"`
}

// ref: http://docs.savla.su/finance
type Item struct {
	ID       	int 	`json:"id"`
	Amount       	float32	`json:"amount"`
	CurrencyCode 	string 	`json:"currency_code"`
	Description 	string 	`json:"description"`
	BillingDate	string 	`json:"billing_date"`
	Misc     	string 	`json:"misc"`
}

// flush finance at start
var finance = Finance{}


func GetAccounts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message": "dumping finance accounts",
		"accounts": finance.Accounts,
	})
}

func GetAccountByOwner(c *gin.Context) {
	owner := c.Param("owner")

	for _, f := range finance.Accounts {
		if f.Owner == owner {
			c.IndentedJSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"account": f,
			})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code": http.StatusNotFound,
		"message": "account not found",
	})
}


func PostDumpRestore(c *gin.Context) {
	var importFinance Finance

	if err := c.BindJSON(&importFinance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	finance = importFinance

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated,
		"message": "finance imported, omitting output",
	})
}

