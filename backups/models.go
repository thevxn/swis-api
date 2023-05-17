package backups

type Backups struct {
	// Array of to be backuped services.
	//Backups []Backup `json:"backups"`
	Backups map[string]Backup `json:"backups"`
}

type Backup struct {
	// Backuped service name -- unique identifier.
	ServiceName string `json:"service_name" binding:"required" validation:"required"`

	// More verbose description of such service backup.
	Description string `json:"description"`

	// Last status string of such backup (e.g. success, failure).
	LastStatus string `json:"last_status" default:"unknown"`

	// UNIX timestamp of the last provided backup.
	Timestamp int `json:"timestamp"`

	// Size of the gzip/tar archive.
	Size string `json:"backup_size"`

	// Name of the compressed backup file.
	FileName string `json:"file_name"`

	// Path to the file on destination machine.
	FileDestination string `json:"file_destination"`

	// Dumping script, git URL.
	ExecutorURL string `json:"executor_url"`

	// Boolean indicating if the service is to be backuped.
	Active bool `json:"active" default:false`
}
