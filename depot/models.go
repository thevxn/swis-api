package depot

type Depots struct {
	Depots []Depot `json:"depots"`
}

type Depot struct {
	Owner      string `json:"owner_name"`
	DepotItems []Item `json:"depot_items"`
}

type Item struct {
	ID          int    `json:"id"`
	Description string `json:"desc"`
	Misc        string `json:"misc"`
	Location    string `json:"depot"`
}

// flush depots at start
var depots = Depots{}
