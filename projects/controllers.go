package projects

import (
	"net/http"

	"go.savla.dev/swis/v5/config"

	"github.com/gin-gonic/gin"
)

var Cache *config.Cache

// @Summary Get all projects
// @Description get project list
// @Tags projects
// @Produce json
// @Success 200 {object} projects.Projects
// @Router /projects [get]
// GetProjects function dumps the projects cache contents.
func GetProjects(c *gin.Context) {
	var projects = Cache.GetAll()

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"message":  "ok, dumping all projects",
		"projects": projects,
	})
	return
}

// @Summary Get project by ID
// @Description get project details by :id route param
// @Tags projects
// @Produce json
// @Success 200 {object} projects.Project
// @Router /projects/{id} [get]
// GetProjectByID returns project's properties, given sent ID exists in database.
func GetProjectByID(c *gin.Context) {
	var id string = c.Param("id")

	rawProject, ok := Cache.Get(id)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "project not found",
			"code":    http.StatusNotFound,
		})
		return
	}

	project, ok := rawProject.(Project)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "cannot assert data type, database internal error",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "dumping requested project's contents",
		"project": project,
	})
	return
}

// @Summary Add new project
// @Description add new project to projects list
// @Tags projects
// @Produce json
// @Param request body projects.Project true "query params"
// @Success 200 {object} projects.Project
// @Router /projects [post]
func PostNewProject(c *gin.Context) {
	var newProject Project

	if err := c.BindJSON(newProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot bind input JSON stream",
		})
		return
	}

	if _, found := Cache.Get(newProject.ID); found {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "project already exists",
			"id":      newProject.ID,
		})
		return
	}

	Cache.Set(newProject.ID, newProject)

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "new project added",
		"project": newProject,
	})
	return
}

// @Summary Update project by its ID
// @Description update project by its ID
// @Tags projects
// @Produce json
// @Param request body projects.Project.ID true "query params"
// @Success 200 {object} projects.Project
// @Router /projects/{id} [put]
func UpdateProjectByID(c *gin.Context) {
	var updatedProject = &Project{}
	var id string = c.Param("id")

	if _, ok := Cache.Get(id); !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "project not found",
			"code":    http.StatusNotFound,
		})
		return
	}

	if err := c.BindJSON(updatedProject); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	Cache.Set(id, updatedProject)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "project updated",
		"project": updatedProject,
	})
	return
}

// @Summary Delete project by its ID
// @Description delete project by its ID
// @Tags projects
// @Produce json
// @Param id path string true "project ID"
// @Success 200 {object} projects.Project.ID
// @Router /projects/{id} [delete]
func DeleteProjectByID(c *gin.Context) {
	var id string = c.Param("id")

	if _, ok := Cache.Get(id); !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "project not found",
			"code":    http.StatusNotFound,
		})
		return
	}

	Cache.Delete(id)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "project deleted by ID",
		"id":      id,
	})
	return
}

// @Summary Upload projects dump -- restore projects
// @Description upload project JSON dump and restore the data model
// @Tags projects
// @Accept json
// @Produce json
// @Router /projects/restore [post]
// PostDumpRestore
func PostDumpRestore(c *gin.Context) {
	var importProjects = &Projects{}
	var project Project
	var counter int = 0

	if err := c.BindJSON(importProjects); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	for _, project = range importProjects.Projects {
		Cache.Set(project.ID, project)
		counter++
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "projects imported/restored, omitting output",
		"count":   counter,
	})
	return
}
