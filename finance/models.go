package finance

type Finance struct {
	Accounts []Account `json:"accounts"`
}

type Account struct {
	Name          string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	// bank codes like "0100" would be invalid as type int
	BankCode string `json:"bank_code"`
	Currency string `json:"account_currency"`
	SWIFT    string `json:"account_swift_bic"`
	IBAN     string `json:"account_iban"`
	Owner    string `json:"account_owner"`
	Items    []Item `json:"account_items"`
}

// ref: http://docs.savla.su/finance
type Item struct {
	ID           int     `json:"id"`
	Amount       float32 `json:"amount"`
	CurrencyCode string  `json:"currency_code"`
	Description  string  `json:"description"`
	BillingDate  string  `json:"billing_date"`
	Misc         string  `json:"misc"`
}

// flush finance at start
var finance = Finance{}
