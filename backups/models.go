package backups

type Backups struct {
	// Array of to be backuped services.
	Backups []Backup `json:"backups"`
}

type Backup struct {
	// Backuped service name -- unique identifier.
	ServiceName string `json:"service_name" binding:"required"`

	// Last status string of such backup (e.g. success, failure).
	LastStatus string `json:"last_status"`

	// UNIX timestamp of the last provided backup.
	Timestamp int `json:"timestamp"`

	// Boolean indicating if the service is to be backuped.
	Active bool `json:"active" default:false`
}

var backups Backups = Backups{}
