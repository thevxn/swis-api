package backups

import (
	"github.com/gin-gonic/gin"
)

func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetBackupStatusAll)
	g.POST("/",
		PostBackedupService)
	g.GET("/:service",
		GetBackedupStatusByServiceName)
	g.PUT("/:service",
		UpdateBackupStatusByServiceName)
	g.DELETE("/:service",
		DeleteBackupByServiceName)
	g.PUT("/:service/active",
		ActiveToggleBackupByServiceName)
	g.POST("/restore",
		PostDumpRestore)
}
