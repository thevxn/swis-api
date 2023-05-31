package system

import "time"

type System struct {
	// Name is a system name.
	Name string `json:"name"`

	// PartCount holds the number of system's parts.
	PartCount int `json:"part_count"`

	// ProjectName is a name of projects.Project object, a pointer to such project.
	ProjectName string `json:"project_name"`

	// Parts is an array of system.Part objects.
	// It is declared as an array to easily index through the items;
	// on the other hand, this makes space for improvement ---
	// --- as searching by name is potentially slower than a direct indexing via maps ---,
	// therefore to be more investigated upon...
	Parts    []Part          `json:"parts"`
	PartsMap map[string]Part `json:"parts_map"`

	// Status is a binary (boolean) status indicator.
	// 1 as OK, 0 as failure.
	Status bool `json:"status_binary"`

	// StatusString is a brief message about the current system state, mode, setting etc.
	StatusString string `json:"status_string"`
}

type Part struct {
	// Name is a part's label.
	Name string `json:"part_name"`

	// Description is a more verbose part commentary.
	Description string `json:"part_description"`
}

//
// v5/system/sync
//

type SyncPacks struct {
	Packs map[string]SyncPack `json:"sync_packs"`
}

type SyncPack struct {
	// Label is the unique name of one sync pack.
	Label string `json:"label"`

	// Timestamp is an UNIX timestamp indicating the last modification datetime.
	// The higher the Timestamp's value, the more possibly the value is taken
	// as the so-called current default val --- successfully winning the race
	// of instances' sync origin per swapi plugin/module.
	//
	// Actually, it would be better to sum sync derivates' values rather then select
	// the newest one as the new base for further timestamps' metadata.
	Timestamp time.Time `json:"timestamp"`

	// Checksums is an array of checksums of defined packages (swapi plugin-modules).
	// See more in the source file /config/sync.go.
	Checksums []string `json:"checksums"`
}
