package roles

type Roles struct {
	//Roles []Role `json:"roles"`
	Roles map[string]Role `json:"roles"`
}

type Role struct {
	// Role name is its unique description, acts like an ID too.
	Name string `json:"name" binding:"required"`

	// Role description to make more sense when listing those.
	Description string `json:"description"`

	// Basic Access-Control List field.
	Admin bool `json:"administrator" default:false`

	// Role status, by default it is inactive.
	Active bool `json:"active" default:false`
}
