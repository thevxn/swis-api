package groups

type Groups struct {
	Groups []Group `json:"groups"`
}

type Group struct {
	ID   string `json:"id"`
	Name string `json:"nickname"`
	Role string `json:"role"`
}

// groups demo data for group struct
var groups = []Group{
	{ID: "1", Name: "superadmins"},
	{ID: "2", Name: "devs"},
	{ID: "3", Name: "ops"},
}
