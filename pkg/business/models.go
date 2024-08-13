package business

// Business structure
type Business struct {
	// Business' unique identifier.
	ID string `json:"id" binding:"required" required:"true" readonly:"true"`

	// Czech Company ICO/ID number.
	ICO int `json:"ico"`

	// Czech company DICO/VAT ID number/string.
	VAT string `json:"vat_id"`

	// Company's business name.
	NameLabel string `json:"name_label"`

	// Array of contacts of different type.
	Contact []Contact `json:"contact"`

	// Business role to such organization (e.g. partner, owner)
	Role string `json:"role"`

	// User's name linked to such business.
	UserName string `json:"username"`
}

type Contact struct {
	// Type of contact field (e.g. e-mail address, street address, telephone number etc).
	Type string `json:"type"`

	// Contact field contents.
	Content string `json:"content"`
}
