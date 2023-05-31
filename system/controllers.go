package system

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// systems (s) and syncs (ss)
var s sync.Map
var ss sync.Map

/*
 *  prototypes
 *
 *  func GetSystems(c *gin.Context) {}
 *  func GetBriefSystemStatus(c *gin.Context) {}
 *  func GetRunningConfiguration(c *gin.Context) {}
 *  func GetSyncTactPackMetadata(c *gin.Context) {}
 *  func CatchSyncTactPack(c *gin.Context) {}
 *  func CatchSyncTactPackByModule(c *gin.Context) {}
 *  func PostNewSyncTactPackByModule(c *gin.Context) {}
 *  func UpdateSyncTactPackByModule(c *gin.Context) {}
 *  func UpdateSyncTactPackByModule(c *gin.Context) {}
 *  func ToggleActiveBoolByModule(c *gin.Context) {}
 *  func DeleteSyncTactPackByModule(c *gin.Context) {}
 *  func PostDumpRestore(c *gin.Context) {}

:6,16s/\(.*\)/func \1(c *gin.Context) {}/
*/

func GetSystems(c *gin.Context) {
	var systems = make(map[string]System)

	s.Range(func(rawKey, rawVal interface{}) bool {
		k, ok := rawKey.(string)
		v, ok := rawVal.(System)

		if !ok {
			return false
		}

		systems[k] = v
		return true
	})

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "dumping systems",
		"systems": systems,
	})
	return
}

func GetBriefSystemStatus(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, papi",
		"status":  0,
	})
	return
}

func GetSyncRunningConfiguration(c *gin.Context) {
	var syncPacks = make(map[string]SyncPack)

	ss.Range(func(rawKey, rawVal interface{}) bool {
		k, ok := rawKey.(string)
		v, ok := rawVal.(SyncPack)

		if !ok {
			return false
		}

		syncPacks[k] = v
		return true
	})

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":       http.StatusOK,
		"message":    "dumping sync packs",
		"sync_packs": syncPacks,
	})
	return
}

func GetSyncTactPackMetadata(c *gin.Context)     {}
func CatchSyncTactPack(c *gin.Context)           {}
func CatchSyncTactPackByModule(c *gin.Context)   {}
func PostNewSyncTactPackByModule(c *gin.Context) {}
func UpdateSyncTactPackByModule(c *gin.Context)  {}
func ToggleActiveBoolByModule(c *gin.Context)    {}
func DeleteSyncTactPackByModule(c *gin.Context)  {}

// @Summary Upload systems dump backup -- restores all systems
// @Description update systems' JSON dump
// @Tags system
// @Accept json
// @Produce json
// @Router /system/restore [post]
// PostDumpRestore
func PostDumpRestoreSystems(c *gin.Context) {
	var importSystems = map[string]System{}
	var system System

	if err := c.BindJSON(importSystems); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	for _, system = range importSystems {
		s.Store(system.Name, system)
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "systems imported/restored, omitting output",
	})
	return
}

func PostDumpRestoreSyncPacks(c *gin.Context) {
	var importSyncPacks = &SyncPacks{
		make(map[string]SyncPack),
	}
	var sync SyncPack

	if err := c.BindJSON(importSyncPacks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	for _, sync = range importSyncPacks.Packs {
		ss.Store(sync.Label, sync)
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "sync packs imported/restored, omitting output",
	})
	return
}
