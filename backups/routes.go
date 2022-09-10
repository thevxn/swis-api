package backups

import (
	"github.com/gin-gonic/gin"
)

func Routes(g *gin.RouterGroup) {
	// @Summary Get all backups status
	// @Description get backups actual status
	// @Tags backup
	// @Produce  json
	// @Success 200 {object} string "ok"
	// @Router /backups [get]
	g.GET("/", GetBackupsStatus)

	// @Summary Get backup status by project/service
	// @Description get backup status by project/service
	// @Router /backups/{service} [get]
	g.GET("/:service", GetBackupStatusByServiceName)

	// @Summary Adding new backuped serivce
	// @Description add new backuped service
	// @Tags backups
	// @Produce json
	// @Param request body backups.Backup true "query params"
	// @Success 200 {object} backups.Backup
	g.POST("/", PostBackupService)

	// @Summary Update backup status by service
	// @Description update backup status by service
	// @Tags backups
	// @Produce json
	// @Param request body backups.Service.Name true "query params"
	// @Success 200 {object} backups.Backup
	// @Router /backups/{service} [put]
	g.PUT("/:service", UpdateBackupStatusByServiceName)
	g.PUT("/:service/active", ActiveToggleBackupByServiceName)

	// @Summary Delete backup service by its Name
	// @Description delete backup service by its Name
	// @Tags backups
	// @Produce json
	// @Success 200 {string} string "ok"
	// @Router /backups/{service} [delete]
	g.DELETE("/:service", DeleteBackupByServiceName)

	// @Summary Upload backups dump backup -- restores all backup services
	// @Description upload backups JSON dump
	// @Tags backups
	// @Accept json
	// @Produce json
	// @Router /backups/restore [post]
	g.POST("/restore", PostDumpRestore)
}
