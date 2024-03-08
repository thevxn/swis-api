package depots

type DepotItem struct {
	// Numeric unique ID of such Item.
	ID int `json:"id" binding:"required"`

	// Item description, name, amount, type etc.
	Description string `json:"desc"`

	// More information, e.g. the more precise location specification.
	Misc string `json:"misc"`

	// Location name of such Item.
	Location string `json:"depot"`

	// Owner name according to users package.
	Owner string `json:"owner_name"`
}

type Location struct {
	// Name of such location, place or depot situation.
	Name string `json:"location_name"`

	// More precise location information.
	Misc string `json:"location_misc"`
}
