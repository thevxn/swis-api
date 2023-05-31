package system

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var s sync.Map

/*
 *  prototypes
 *
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

func GetBriefSystemStatus(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok, papi",
		"status":  0,
	})
	return
}

func GetRunningConfiguration(c *gin.Context) {
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

func GetSyncTactPackMetadata(c *gin.Context)     {}
func CatchSyncTactPack(c *gin.Context)           {}
func CatchSyncTactPackByModule(c *gin.Context)   {}
func PostNewSyncTactPackByModule(c *gin.Context) {}
func UpdateSyncTactPackByModule(c *gin.Context)  {}
func ToggleActiveBoolByModule(c *gin.Context)    {}
func DeleteSyncTactPackByModule(c *gin.Context)  {}
func PostDumpRestore(c *gin.Context)             {}
