package projects

type Project struct {
	// Project's unique identificator.
	ID string `json:"project_id" binding:"required" validation:"required" required:"true"`

	// Project name.
	Name string `json:"project_name" binding:"required" validation:"required" required:"true"`

	// Brief project description.
	Description string `json:"project_desc"`

	// URL to documentation page(s).
	DocsLink string `json:"project_docs_link"`

	// Project manager's name/username.
	Manager string `json:"project_manager"`

	// Published boolean.
	Published bool `json:"project_published" default:false`

	// Git repository link (not URL, without HTTP scheme).
	Repository string `json:"project_repo"`

	// URL to redmine project overview.
	Redmine string `json:"redmine_link"`

	// URL to kanboard/kanban project's page.
	Kanban string `json:"kanban_link"`

	// Projects backuped boolean.
	Backuped bool `json:"backuped" default:false`

	// URL to base page of the project (project's URL).
	URL string `json:"project_url"`

	// Target internal node of deployment.
	DeployTarget string `json:"project_deploy_target"`
}
