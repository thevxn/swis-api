package projects

type Projects struct {
	Projects []Project `json:"projects"`
}

type Project struct {
	ID           string `json:"project_id"`
	Name         string `json:"project_name"`
	Description  string `json:"project_desc"`
	DocsLink     string `json:"project_docs_link"`
	Manager      string `json:"project_manager"`
	Published    bool   `json:"project_published" default:false`
	Repository   string `json:"project_repo"`
	URL          string `json:"project_url"`
	DeployTarget string `json:"project_deploy_target"`
}

// flush projects object/array
var projects = Projects{}
