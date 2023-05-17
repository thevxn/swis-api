package depots

type Depots struct {
	Depots map[string]Depot `json:"depots"`
}

type Depot struct {
	// Depot owner's name, unique ID.
	Owner string `json:"owner_name"`

	// Generic array of depot Items.
	DepotItems []Item `json:"depot_items"`
}

type Item struct {
	// Numeric unique ID of such Item.
	ID int `json:"id"`

	// Item description, name, amount, type etc.
	Description string `json:"desc"`

	// More information, e.g. the more precise location specification.
	Misc string `json:"misc"`

	// Location name of such Item.
	Location string `json:"depot"`
}

type Location struct {
	// Name of such location, place or depot situation.
	Name string `json:"location_name"`

	// More precise location information.
	Misc string `json:"location_misc"`
}
