package backups

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get all backups status
// @Description get backups actual status
// @Tags backup
// @Produce  json
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
func GetBackupStatusByService(c *gin.Context) {}

// @Summary Adding new backuped serivce
// @Description add new backuped service
// @Tags backups
// @Produce json
// @Param request body backups.Backup true "query params"
// @Success 200 {object} backups.Backup
func PostBackupService(c *gin.Context) {}

// @Summary Update backup status by service
// @Description update backup status by service
// @Tags backups
// @Produce json
// @Param request body backups.Backup.ServiceName true "query params"
// @Success 200 {object} backups.Backup
// @Router /backups/{service} [put]
func UpdateBackupStatusByService(c *gin.Context) {}

// @Summary Delete backup service by its Name
// @Description delete backup service by its Name
// @Tags backups
// @Produce json
// @Success 200 {string} string "ok"
// @Router /backups/{service} [delete]
func DeleteBackupByService(c *gin.Context) {}

// @Summary Upload backups dump backup -- restores all backup services
// @Description upload backups JSON dump
// @Tags backups
// @Accept json
// @Produce json
// @Router /backups/restore [post]
func PostDumpRestore(c *gin.Context) {}
