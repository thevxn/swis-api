package finance

type Finance struct {
	// Finance accounts list.
	Accounts []Account `json:"accounts"`
}

type Account struct {
	// Account name/unique ID.
	Name string `json:"account_name"`

	// Account number.
	AccountNumber string `json:"account_number"`

	// Universal in-state bank code (CZ mainly).
	// Bank codes such as "0100" would be invalid as type int!
	BankCode string `json:"bank_code"`

	// Account currency name (e.g. CZK, GBP, EUR, USD)
	Currency string `json:"account_currency"`

	// Account SWIFT/BIC code for international payments.
	SWIFT string `json:"account_swift_bic"`

	// Account IBAN code for international payments.
	IBAN string `json:"account_iban"`

	// Owner's name/username to link account to.
	Owner string `json:"account_owner"`

	// Account items like (incoming/outcoming) payments.
	Items []Item `json:"account_items"`
}

// ref: http://docs.savla.su/finance
type Item struct {
	// Item unique ID (e.g. datetime timestamp plus currency etc).
	ID int `json:"id"`

	// Payment amount in defined currency (often the account's currency).
	Amount float32 `json:"amount"`

	// Payment currency name (e.g. CZK, GBP, EUR, USD)
	Currency string `json:"currency"`

	// Payment/item description.
	Description string `json:"description"`

	// Billing date of the payment.
	BillingDate string `json:"billing_date"`

	// Mescellaneous information about the payment (e.g. foreign currency and amount).
	Misc string `json:"misc"`
}

// flush finance at start
var finance = Finance{}
