package business

type BusinessArray struct {
	BusinessArray []Business `json:"business"`
}

// Business structure
type Business struct {
	ID        string    `json:"id"`
	ICO       string    `json:"ico"`
	VAT       string    `json:"vat_id"`
	NameLabel string    `json:"name_label"`
	Contact   []Contact `json:"contact"`
	Role      string    `json:"role"`
	UserName  string    `json:"username"`
	//TokenBase64		string 	`json:"token_base64"`
}

type Contact struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// flush businessArray at start
var businessArray = BusinessArray{}
