// internet of things (iot) module + home module.
package iot

type Things struct {
	// IoTs' owner name.
	Owner string `json:"owner_name"`

	// Array/map of Things.
	Things []Thing `json:"things"`
	//Things map[string]Thing `json:"things"`
}

type Thing struct {
	// Unique thing's hash ID.
	Hash string `json:"thing_hash" binding:"required"`

	// Thing's name, descriptor.
	Name string `json:"thing_name"`

	// Thing's more verbose description.
	Description string `json:"thing_desc"`

	// Type of thing (e.g. sensor)
	Type string `json:"thing_type"`

	// Target bus device (e.g. raspi01.savla.iot)
	Bus string `json:"thing_bus"`
}

type Bus struct{}
