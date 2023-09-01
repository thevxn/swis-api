package system

import (
	"github.com/gin-gonic/gin"
)

func Routes(g *gin.RouterGroup) {
	g.GET("/",
		GetSystems)
	g.GET("/status",
		GetBriefSystemStatus)
	g.GET("/sync",
		GetSyncRunningConfiguration)
	g.GET("/sync/:module",
		GetSyncTactPackMetadata)
	g.PUT("/sync",
		CatchSyncTactPack)
	g.PUT("/sync/:module",
		CatchSyncTactPackByModule)
	g.POST("/sync/:module",
		PostNewSyncTactPackByModule)
	//g.PUT("/sync/:module",
	//	UpdateSyncTactPackByModule)
	g.PATCH("/sync/:module",
		UpdateSyncTactPackByModule)
	g.PUT("/sync/:module/active",
		ToggleActiveBoolByModule)
	g.DELETE("/sync/:module",
		DeleteSyncTactPackByModule)
	g.POST("/restore",
		PostDumpRestoreSystems)
	g.POST("/sync/restore",
		PostDumpRestoreSyncPacks)
}
