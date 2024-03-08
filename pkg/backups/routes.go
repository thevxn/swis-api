package backups

import (
	"github.com/gin-gonic/gin"
)

func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetBackups)
	g.POST("/:key",
		PostBackedupServiceByServiceKey)
	g.GET("/:key",
		GetBackedupStatusByServiceKey)
	g.PUT("/:key",
		UpdateBackupStatusByServiceKey)
	g.DELETE("/:key",
		DeleteBackupByServiceKey)
	g.PUT("/:key/active",
		ActiveToggleBackupByServiceKey)
	g.POST("/restore",
		PostDumpRestore)
}
