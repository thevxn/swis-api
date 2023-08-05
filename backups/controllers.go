package backups

import (
	"net/http"

	"go.savla.dev/swis/v5/config"

	"github.com/gin-gonic/gin"
)

var Cache *config.Cache

// @Summary Get all backed up services
// @Description get backed up services
// @Tags backups
// @Produce json
// @Success 200 {object} string "ok"
// @Router /backups [get]
func GetBackupStatusAll(c *gin.Context) {
	backedupServices, count := Cache.GetAll()

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"count":   count,
		"message": "ok, dumping all backed up services",
		"backups": backedupServices,
	})
	return
}

// @Summary Get backup status by project/service
// @Description get backup status by project/service
// @Tags backups
// @Produce  json
// @Param   host     path    string     true        "backup service name"
// @Success 200 {string} string	"ok"
// @Router /backups/{service} [get]
func GetBackedupStatusByServiceName(c *gin.Context) {
	var name string = c.Param("service")
	var backedupService Backup

	rawService, found := Cache.Get(name)
	if !found {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "backup status by service not found",
		})
		return
	}

	backedupService, ok := rawService.(Backup)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "cannot assert data type, database internal error",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "dumping requested backed up service's status",
		"backup":  backedupService,
	})
	return
}

// @Summary Add new backed up serivce
// @Description add new backed up service
// @Tags backups
// @Produce json
// @Param request body backups.Backup true "query params"
// @Success 200 {object} backups.Backup
// @Router /backups [post]
func PostBackedupService(c *gin.Context) {
	var newBackedupService Backup

	if err := c.BindJSON(&newBackedupService); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	if _, found := Cache.Get(newBackedupService.ServiceName); found {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "backed up service already exists",
			"name":    newBackedupService.ServiceName,
		})
		return
	}

	if saved := Cache.Set(newBackedupService.ServiceName, newBackedupService); !saved {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "backed up service couldn't be saved to database",
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "new backed up service added",
		"service": newBackedupService,
	})
	return
}

// @Summary Update backup status by service
// @Description update backup status by service
// @Tags backups
// @Produce json
// @Param request body backups.Backup.ServiceName true "query params"
// @Success 200 {object} backups.Backup
// @Router /backups/{service} [put]
func UpdateBackupStatusByServiceName(c *gin.Context) {
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
func ActiveToggleBackupByServiceName(c *gin.Context) {
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

// @Summary Delete backup service by its Name
// @Description delete backup service by its Name
// @Tags backups
// @Produce json
// @Success 200 {string} string "ok"
// @Router /backups/{service} [delete]
func DeleteBackupByServiceName(c *gin.Context) {
	var name string = c.Param("service")

	if _, found := Cache.Get(name); !found {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "backed up service not found",
			"name":    name,
		})
		return
	}

	if deleted := Cache.Delete(name); !deleted {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "backed up service couldn't be deleted from database",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "backed up service deleted by Name",
		"name":    name,
	})
	return
}

// @Summary Upload backups dump backup -- restores all backup services
// @Description upload backups JSON dump
// @Tags backups
// @Accept json
// @Produce json
// @Router /backups/restore [post]
func PostDumpRestore(c *gin.Context) {
	var importServices = &Backups{}
	var service Backup

	if err := c.BindJSON(importServices); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	for _, service = range importServices.Backups {
		Cache.Set(service.ServiceName, service)
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "backuped services imported, omitting output",
	})
	return
}
