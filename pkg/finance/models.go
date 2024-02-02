package finance

import (
	"time"
)

type Account struct {
	// Account unique ID, snake_cased identifier.
	ID string `json:"account_id" binding:"required"`

	// Account more verbose name.
	Description string `json:"account_description"`

	// Account number.
	AccountNumber string `json:"account_number" binding:"required"`

	// Universal in-state bank code (CZ mainly).
	// Bank codes such as "0100" would be invalid as type int!
	BankCode string `json:"bank_code"`

	// Account currency name (e.g. CZK, GBP, EUR, USD)
	Currency string `json:"account_currency" binidng:"required"`

	// Account SWIFT/BIC code for international payments.
	SWIFT string `json:"account_swift_bic" binding:"required"`

	// Account IBAN code for international payments.
	IBAN string `json:"account_iban" binding:"required"`

	// Owner's name/username to link account to.
	Owner string `json:"account_owner"`
}

// ref: http://docs.savla.su/finance
type Item struct {
	// Item unique ID (e.g. datetime timestamp plus currency etc).
	ID string `json:"id" binding:"required"`

	// Type defines whether the item is an income, or an expense.
	Type string `json:"type" binding:"required"`

	// Payment amount in defined currency (often the account's currency).
	Amount float64 `json:"amount" binding:"required"`

	// Payment currency name (e.g. CZK, GBP, EUR, USD).
	Currency string `json:"currency" binding:"required"`

	// Payment/item description.
	Description string `json:"description"`

	// PaymentDate is a datetime of the payment.
	PaymentDate time.Time `json:"payment_date" binding:"required"`

	// Referencing finance account.
	AccountID string `json:"account_id" binding:"required"`

	// BusinessID is a reference to 'business' package.
	BusinessID string `json:"business_id"`

	// Invoice identification.
	InvoiceNo string `json:"invoice_no"`

	// Mescellaneous information about the payment (e.g. foreign currency and amount).
	Misc string `json:"misc"`
}

type Tax struct {
	// Sum of incomes.
	IncomeTotal float64 `json:"income_total"`

	// Sum of expenses.
	ExpenseTotal float64 `json:"expense_total"`

	// Difference between incomes and expenses. Base for income tax.
	Base float64 `json:"base_sum"`
}
