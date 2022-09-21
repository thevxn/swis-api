package links

type Links struct {
	// An array of links from JSON.
	Links []Link `json:"links"`
}

type Link struct {
	// Link's unique hash ID.
	Hash string `json:"hash" binding:"required" validation:"required"`

	// Link's name/placeholder.
	Name string `json:"name"`

	// Link's more verbose name/description.
	Description string `json:"description"`

	// Link's URL to link to.
	URL string `json:"url" binding:"required"`

	// Link's activated status.
	Active bool `json:"active" default:false`
}

var links = []Link{}
