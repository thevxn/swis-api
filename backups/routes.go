package backups

import (
	"github.com/gin-gonic/gin"
)

func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetBackupsStatus)
	g.POST("/",
		PostBackupService)
	g.GET("/:service",
		GetBackupStatusByServiceName)
	g.PUT("/:service",
		UpdateBackupStatusByServiceName)
	g.DELETE("/:service",
		DeleteBackupByServiceName)
	g.PUT("/:service/active",
		ActiveToggleBackupByServiceName)
	g.POST("/restore",
		PostDumpRestore)
}
