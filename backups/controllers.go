package backups

import (
	"net/http"

	"go.savla.dev/swis/v5/config"

	"github.com/gin-gonic/gin"
)

var (
	Cache   *config.Cache
	pkgName string = "backups"
)

// @Summary Get all backed up services
// @Description get backed up services
// @Tags backups
// @Produce json
// @Success 200 {object} string "ok"
// @Router /backups [get]
func GetBackupStatusAll(ctx *gin.Context) {
	config.PrintAllRootItems(ctx, Cache, pkgName)
	return
}

// @Summary Get backup status by project's/service's key
// @Description get backup status by project'S/service's key
// @Tags backups
// @Produce json
// @Param host path string true "backup service key"
// @Success 200 {string} string	"ok"
// @Router /backups/{key} [get]
func GetBackedupStatusByServiceKey(ctx *gin.Context) {
	config.PrintItemByParam(ctx, Cache, pkgName, Backup{})
	return
}

// @Summary Add new backed up serivce
// @Description add new backed up service
// @Tags backups
// @Produce json
// @Param request body backups.Backup true "query params"
// @Success 200 {object} backups.Backup
// @Router /backups/{key} [post]
func PostBackedupServiceByServiceKey(ctx *gin.Context) {
	config.AddNewItemByParam(ctx, Cache, pkgName, Backup{})
	return
}

// @Summary Delete backup service by its key
// @Description delete backup service by its key
// @Tags backups
// @Produce json
// @Success 200 {string} string "ok"
// @Router /backups/{key} [delete]
func DeleteBackupByServiceKey(ctx *gin.Context) {
	config.DeleteItemByParam(ctx, Cache, pkgName)
	return
}

// @Summary Upload backups dump backup -- restores all backup services
// @Description upload backups JSON dump
// @Tags backups
// @Accept json
// @Produce json
// @Router /backups/restore [post]
func PostDumpRestore(ctx *gin.Context) {
	config.BatchRestoreItems(ctx, Cache, pkgName, Backup{})
	return
}

// @Summary Update backup status by service's key
// @Description update backup status by service's key
// @Tags backups
// @Produce json
// @Param request body backups.Backup.ServiceName true "query params"
// @Success 200 {object} backups.Backup
// @Router /backups/{key} [put]
func UpdateBackupStatusByServiceKey(c *gin.Context) {
	var updatedService Backup
	var postedService Backup

	var name string = c.Param("service")

	rawService, found := Cache.Get(name)
	if !found {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "backed up service not found by its name",
			"name":    name,
		})
		return
	}

	updatedService, ok := rawService.(Backup)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "cannot assert data type, database internal error",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	if err := c.BindJSON(&postedService); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	// manually update important report fields
	// TODO: review this!
	updatedService.Timestamp = postedService.Timestamp
	updatedService.LastStatus = postedService.LastStatus
	updatedService.FileName = postedService.FileName
	updatedService.Size = postedService.Size

	if saved := Cache.Set(name, updatedService); !saved {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "backed up service couldn't be saved to database",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "backed up service updated",
		"backup":  updatedService,
	})
	return
}

// (PUT /backups/{service}/active)
// @Summary Acitive/inactive backup toggle by its ServiceName
// @Description active/inactive backup toggle by its ServiceName
// @Tags backups
// @Produce json
// @Param  service_name  path  string  true  "service name"
// @Success 200 {object} backups.Backup
// @Router /backups/{service}/active [put]
func ActiveToggleBackupByServiceKey(c *gin.Context) {
	var service Backup
	var name string = c.Param("service")

	rawService, found := Cache.Get(name)
	if !found {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "backed up service not found",
			"code":    http.StatusNotFound,
			"name":    name,
		})
		return
	}

	service, ok := rawService.(Backup)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "cannot assert data type, database internal error",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// inverse the Active field value
	service.Active = !service.Active

	if saved := Cache.Set(name, service); !saved {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "backed up service couldn't be saved to database",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "backed up service active toggle pressed!",
		"backup":  service,
	})
	return
}
