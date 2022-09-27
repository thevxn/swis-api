package projects

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
		"code":    http.StatusNotFound,
		"message": "project not found",
	})
	return nil
}

// @Summary Get all projects
// @Description get project list
// @Tags projects
// @Produce  json
// @Success 200 {object} projects.Projects
// @Router /projects/{name} [get]
// GetProjects function dumps the projects variable contents
func GetProjects(c *gin.Context) {
	// serialize struct to JSON
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"message":  "dumping projects",
		"projects": projects.Projects,
	})
}

// @Summary Get project by ID
// @Description get project details by :id param
// @Tags projects
// @Produce  json
// @Success 200 {object} projects.Project
// @Router /projects/{id} [get]
// GetProjectByID returns project's properties, given sent ID exists in database
func GetProjectByID(c *gin.Context) {
	if p := findProjectByID(c); p != nil {
		// project found
		c.IndentedJSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "dumping requested project's contents",
			"project": p,
		})
	}
}

// @Summary Add new project
// @Description add new project to projects list
// @Tags projects
// @Produce json
// @Param request body projects.Project true "query params"
// @Success 200 {object} projects.Project
// @Router /projects [post]
// PostProject
func PostProject(c *gin.Context) {
	var newProject Project

	// bind received JSON to newProject
	if err := c.BindJSON(&newProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	// add new project
	projects.Projects = append(projects.Projects, newProject)
	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "new project added",
		"project": newProject,
	})
}

// @Summary Upload projects dump -- restore projects
// @Description upload project JSON dump and restore the data model
// @Tags projects
// @Accept json
// @Produce json
// @Router /projects/restore [post]
// PostDumpRestore
func PostDumpRestore(c *gin.Context) {
	var importProjects Projects

	if err := c.BindJSON(&importProjects); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	projects = importProjects

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "projects imported/restored, omitting output",
	})
}
