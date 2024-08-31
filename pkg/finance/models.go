package finance

import (
	"time"
)

type Account struct {
	// Account unique ID, snake_cased identifier.
	ID string `json:"id" binding:"required" required:"true" readonly:"true"`

	// Account more verbose name.
	Description string `json:"account_description"`

	// Account number.
	AccountNumber string `json:"account_number" binding:"required" required:"true"`

	// Universal in-state bank code (CZ mainly).
	// Bank codes such as "0100" would be invalid as type int!
	BankCode string `json:"bank_code"`

	// Account currency name (e.g. CZK, GBP, EUR, USD)
	Currency string `json:"account_currency" binidng:"required" required:"true"`

	// Account SWIFT/BIC code for international payments.
	SWIFT string `json:"account_swift_bic" binding:"required" required:"true"`

	// Account IBAN code for international payments.
	IBAN string `json:"account_iban" binding:"required" required:"true"`

	// Owner's name/username to link account to.
	Owner string `json:"account_owner"`
}

// ref: http://docs.vxn.su/finance
type Item struct {
	// Item unique ID (e.g. datetime timestamp plus currency etc).
	ID string `json:"id" binding:"required" required:"true" readonly:"true"`

	// Type defines whether the item is an income, or an expense.
	Type string `json:"type" binding:"required" required:"true"`

	// Payment amount in defined currency (often the account's currency).
	Amount float64 `json:"amount" binding:"required" required:"true"`

	// Payment currency name (e.g. CZK, GBP, EUR, USD).
	Currency string `json:"currency" binding:"required" required:"true"`

	// Payment/item description.
	Description string `json:"description"`

	// PaymentDate is a datetime of the payment.
	PaymentDate time.Time `json:"payment_date" binding:"required" required:"true"`

	// Referencing finance account.
	AccountID string `json:"account_id" binding:"required" required:"true"`

	// BusinessID is a reference to 'business' package.
	BusinessID string `json:"business_id"`

	// Invoice identification.
	InvoiceNo string `json:"invoice_no"`

	// Tags to filter items.
	Tags []string `json:"tags"`

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

	// Absolute estimation of the income tax per year, from given totals.
	IncomeTaxEstAbs float64 `json:"income_tax_estimation_abs"`

	// Weighted estimation of the income tax, where 60% of income are meant as expenses.
	IncomeTaxEst60 float64 `json:"income_tax_estimation_60"`
}
