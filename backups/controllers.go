package backups

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var b sync.Map

// @Summary Get all backuped services
// @Description get backuped services
// @Tags backups
// @Produce json
// @Success 200 {object} string "ok"
// @Router /backups [get]
func GetBackupsStatus(c *gin.Context) {
	var services = make(map[string]Backup)

	b.Range(func(rawKey, rawVal interface{}) bool {
		k, ok := rawKey.(string)
		v, ok := rawVal.(Backup)

		if !ok {
			return false
		}

		services[k] = v
		return true
	})

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, dumping all sockets",
		"backups": services,
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
func GetBackupStatusByServiceName(c *gin.Context) {
	var name string = c.Param("service")
	var service Backup

	rawService, ok := b.Load(name)
	service, ok = rawService.(Backup)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "service not found",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "dumping requested backuped service's status",
		"backup":  service,
	})
	return
}

// @Summary Adding new backuped serivce
// @Description add new backuped service
// @Tags backups
// @Produce json
// @Param request body backups.Backup true "query params"
// @Success 200 {object} backups.Backup
// @Router /backups [post]
func PostBackupService(c *gin.Context) {
	var newService = &Backup{}

	if err := c.BindJSON(newService); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	if _, found := b.Load(newService.ServiceName); found {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "service already exists",
			"name":    newService.ServiceName,
		})
		return
	}

	b.Store(newService.ServiceName, newService)

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "new project added",
		"service": newService,
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
	var postedService = &Backup{}
	var name string = c.Param("service")

	rawService, ok := b.Load(name)
	updatedService, ok = rawService.(Backup)

	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "service not found",
			"name":    name,
		})
		return
	}

	if err := c.BindJSON(postedService); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	// manually update important report fields
	updatedService.Timestamp = postedService.Timestamp
	updatedService.LastStatus = postedService.LastStatus
	updatedService.FileName = postedService.FileName

	b.Store(name, updatedService)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "service updated",
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

	rawService, ok := b.Load(name)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "service not found",
			"code":    http.StatusNotFound,
			"name":    name,
		})
		return
	}

	service, typeOk := rawService.(Backup)
	if !typeOk {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"message": "stored value is not type Backup",
			"code":    http.StatusConflict,
		})
		return
	}

	// inverse the Active field value
	service.Active = !service.Active

	b.Store(name, service)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "backuped service active toggle pressed!",
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

	if _, ok := b.Load(name); !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "service not found",
			"name":    name,
		})
		return
	}

	b.Delete(name)

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "service deleted by Name",
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
		b.Store(service.ServiceName, service)
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "backuped services imported, omitting output",
	})
	return
}
