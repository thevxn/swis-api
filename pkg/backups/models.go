package backups

type Backup struct {
	// Unique identifier.
	ID string `json:"id" binding:"required" validation:"required" required:"true" readonly:"true"`

	// Backuped service name -- unique identifier.
	ServiceName string `json:"service_name" binding:"required" validation:"required" required:"true"`

	// More verbose description of such service backup.
	Description string `json:"description"`

	// Last status string of such backup (e.g. success, failure).
	LastStatus string `json:"last_status" default:"unknown"`

	// UNIX timestamp of the last provided backup.
	Timestamp int `json:"timestamp"`

	// Size of the gzip/tar archive in bytes.
	Size int `json:"backup_size"`

	// Name of the compressed backup file.
	FileName string `json:"file_name"`

	// Path to the file on destination machine.
	FileDestination string `json:"file_destination"`

	// Dumping script, git URL.
	ExecutorURL string `json:"executor_url"`

	// Reference to projects pkg's instance.
	ProjectID string `json:"project_id"`

	// Boolean indicating if the service is to be backuped.
	Active bool `json:"active" default:false`
}
