package projects

import (
	"go.savla.dev/swis/v5/config"

	"github.com/gin-gonic/gin"
)

var (
	Cache   *config.Cache
	pkgName string = "projects"
)

// GetProjects function dumps the projects cache contents.
// @Summary Get all projects
// @Description get project list
// @Tags projects
// @Produce json
// @Success 200 {object} []projects.Project
// @Router /projects [get]
func GetProjects(ctx *gin.Context) {
	config.PrintAllRootItems(ctx, Cache, pkgName)
	return
}

// GetProjectByKey returns project's properties, given sent ID exists in database.
// @Summary Get project by ID
// @Description get project details by :id route param
// @Tags projects
// @Produce json
// @Success 200 {object} projects.Project
// @Router /projects/{key} [get]
func GetProjectByKey(ctx *gin.Context) {
	config.PrintItemByParam(ctx, Cache, pkgName, Project{})
	return
}

// @Summary Add new project
// @Description add new project to projects list
// @Tags projects
// @Produce json
// @Param request body projects.Project true "query params"
// @Success 200 {object} projects.Project
// @Router /projects/{key} [post]
func PostNewProjectByKey(ctx *gin.Context) {
	config.AddNewItemByParam(ctx, Cache, pkgName, Project{})
	return
}

// @Summary Update project by its ID
// @Description update project by its ID
// @Tags projects
// @Produce json
// @Param request body projects.Project.ID true "query params"
// @Success 200 {object} projects.Project
// @Router /projects/{key} [put]
func UpdateProjectByKey(ctx *gin.Context) {
	config.UpdateItemByParam(ctx, Cache, pkgName, Project{})
	return
}

// @Summary Delete project by its ID
// @Description delete project by its ID
// @Tags projects
// @Produce json
// @Param id path string true "project ID"
// @Success 200 {object} projects.Project.ID
// @Router /projects/{key} [delete]
func DeleteProjectByKey(ctx *gin.Context) {
	config.DeleteItemByParam(ctx, Cache, pkgName)
	return
}

// PostDumpRestore
// @Summary Upload projects dump -- restore projects
// @Description upload project JSON dump and restore the data model
// @Tags projects
// @Accept json
// @Produce json
// @Router /projects/restore [post]
func PostDumpRestore(ctx *gin.Context) {
	config.BatchRestoreItems(ctx, Cache, pkgName, Project{})
	return
}
