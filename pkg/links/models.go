package links

type Link struct {
	// Link's name/placeholder, used as an unique key (legacy).
	ID string `json:"id" binding:"required" required:"true" readonly:"true"`

	// Link's name/placeholder, used as an unique key (legacy).
	Name string `json:"name" binding:"required" required:"true"`

	// Link's more verbose name/description.
	Description string `json:"description"`

	// Link's URL to link to.
	URL string `json:"url" binding:"required" required:"true"`

	// Link's activated status.
	Active bool `json:"active" default:false`
}
