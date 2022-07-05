package projects

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Projects struct {
	Projects	[]Project	`json:"projects"`
}

type Project struct {
	ID		string	`json:"project_id"`
	Name		string	`json:"project_name"`
	Description	string	`json:"project_desc"`
        DocsLink	string	`json:"project_docs_link"`
        Manager		string	`json:"project_manager"`
        Published	bool	`json:"project_published" default:false`
	Repository	string	`json:"project_repo"`
	URL		string	`json:"project_url"`
	DeployTarget	string	`json:"project_deploy_target"`
}

// flush projects object/array
var projects = Projects{}


func findProjectByID(c *gin.Context) (p *Project) {
	// loop over projects
	for _, p := range projects.Projects {
		if p.ID == c.Param("id") {
			//c.IndentedJSON(http.StatusOK, p)
			return &p
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code": http.StatusNotFound,
		"message": "project not found",
	})
	return nil
}

// GetProjects function dumps the projects variable contents
func GetProjects(c *gin.Context) {
	// serialize struct to JSON
	c.IndentedJSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message": "dumping projects",
		"projects": projects.Projects,
	})
}

// GetProjectByID returns project's properties, given sent ID exists in database
func GetProjectByID(c *gin.Context) {
	if p := findProjectByID(c); p != nil {
		// project found
		c.IndentedJSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"message": "dumping requested project's contents",
			"project": p,
		})
	}
}

// PostProject
func PostProject(c *gin.Context) {
	var newProject Project

	// bind received JSON to newProject
	if err := c.BindJSON(&newProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	// add new project
	projects.Projects = append(projects.Projects, newProject)
	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated,
		"message": "new project added",
		"project": newProject,
	})
}

// PostDumpRestore
func PostDumpRestore(c *gin.Context) {
	var importProjects Projects

	if err := c.BindJSON(&importProjects); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	projects = importProjects

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated,
		"message": "projects imported/restored, omitting output",
	})
}

