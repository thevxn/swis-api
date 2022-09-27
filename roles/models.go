package roles

type Roles struct {
	Roles []Role `json:"roles"`
}

type Role struct {
	// Role name is its unique description, acts like an ID too.
	Name string `json:"name"`

	// Role description to make more sense when listing those.
	Description string `json:"description"`

	// Role status, by default it is inactive.
	Active bool `json:"active" default:false`
}
