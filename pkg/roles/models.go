package roles

type Role struct {
	// Role ID is its unique identifier.
	ID string `json:"name" binding:"required" required:"true" readonly:"true"`

	// Role name is its unique description, acts like an ID too (legacy).
	Name string `json:"name" binding:"required" required:"true"`

	// Role description to make more sense when listing those.
	Description string `json:"description"`

	// Basic Access-Control List field.
	Admin bool `json:"administrator" default:false`

	// Role status, by default it is inactive.
	Active bool `json:"active" default:false`
}
