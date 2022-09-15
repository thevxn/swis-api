package backups

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func findBackupByServiceName(c *gin.Context) (index *int, backup *Backup) {
	for i, b := range backups.Backups {
		if b.ServiceName == c.Param("service") {
			return &i, &b
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "backuped service not found",
	})
	return nil, nil
}

// @Summary Get all backups status
// @Description get backups actual status
// @Tags backups
// @Produce json
// @Success 200 {object} string "ok"
// @Router /backups [get]
func GetBackupsStatus(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, dumping all backuped services",
		"backups": backups.Backups,
	})
}

// @Summary Get backup status by project/service
// @Description get backup status by project/service
// @Tags backups
// @Produce  json
// @Param   host     path    string     true        "dish instance name"
// @Success 200 {string} string	"ok"
// @Router /backups/status/{service} [get]
func GetBackupStatusByServiceName(c *gin.Context) {
	_, backup := findBackupByServiceName(c.Copy())
	if backup == nil {
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, returning found backup status",
		"backup":  *backup,
	})
}

// @Summary Adding new backuped serivce
// @Description add new backuped service
// @Tags backups
// @Produce json
// @Param request body backups.Backup true "query params"
// @Success 200 {object} backups.Backup
func PostBackupService(c *gin.Context) {
	var newBackup Backup

	// bind JSON to newSocket
	if err := c.BindJSON(&newBackup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	// add new backup
	backups.Backups = append(backups.Backups, newBackup)

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "backup added",
		"backup":  newBackup,
	})
}

// @Summary Update backup status by service
// @Description update backup status by service
// @Tags backups
// @Produce json
// @Param request body backups.Backup.ServiceName true "query params"
// @Success 200 {object} backups.Backup
// @Router /backups/{service} [put]
func UpdateBackupStatusByServiceName(c *gin.Context) {
	var updatedBackup Backup

	i, b := findBackupByServiceName(c.Copy())

	if err := c.BindJSON(&updatedBackup); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	if !b.Active {
		c.IndentedJSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "this service is set to inactive",
		})
		return
	}

	// update only some fields, not everything!
	backups.Backups[*i].LastStatus = updatedBackup.LastStatus
	backups.Backups[*i].Timestamp = updatedBackup.Timestamp
	backups.Backups[*i].Size = updatedBackup.Size
	backups.Backups[*i].FileName = updatedBackup.FileName

	c.IndentedJSON(http.StatusAccepted, gin.H{
		"code":    http.StatusAccepted,
		"message": "backup status updated",
		"backup":  backups.Backups[*i],
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
	var updatedBackup Backup

	i, _ := findBackupByServiceName(c.Copy())
	updatedBackup = backups.Backups[*i]

	// inverse the Muted field value
	updatedBackup.Active = !updatedBackup.Active

	backups.Backups[*i] = updatedBackup
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "backup service activated toggle pressed!",
		"backup":  updatedBackup,
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
	i, b := findBackupByServiceName(c.Copy())

	// delete an element from the array
	// https://www.educative.io/answers/how-to-delete-an-element-from-an-array-in-golang
	newLength := 0
	for index := range backups.Backups {
		if *i != index {
			backups.Backups[newLength] = backups.Backups[index]
			newLength++
		}
	}

	// reslice the array to remove extra index
	backups.Backups = backups.Backups[:newLength]

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "backup deleted by ServiceName",
		"backup":  *b,
	})
}

// @Summary Upload backups dump backup -- restores all backup services
// @Description upload backups JSON dump
// @Tags backups
// @Accept json
// @Produce json
// @Router /backups/restore [post]
func PostDumpRestore(c *gin.Context) {
	var importBackups Backups

	// bind received JSON to importBackups
	if err := c.BindJSON(&importBackups); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	backups = importBackups

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "backups imported successfully",
	})
}
