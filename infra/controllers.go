package infra

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func findHostByHostname(c *gin.Context) (index *int, h *Host) {
	// loop over hosts
	var hosts = infrastructure.Hosts

	for i, a := range hosts {
		if a.HostnameShort == c.Param("hostname") || a.HostnameFQDN == c.Param("hostname") {
			//c.IndentedJSON(http.StatusOK, a)
			return &i, &a
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "host not found",
	})
	return nil, nil
}

func GetInfrastructure(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"infrastructure": infrastructure,
	})
}

func GetHosts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"hosts": infrastructure.Hosts,
	})
}

func GetNetworks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"networks": infrastructure.Networks,
	})
}

func GetHostByHostname(c *gin.Context) {
	if _, host := findHostByHostname(c); host != nil {
		// host found
		c.IndentedJSON(http.StatusOK, host)
	}
}

// PostDumpRestore
func PostDumpRestore(c *gin.Context) {
	var importInfrastructure Infrastructures

	if err := c.BindJSON(&importInfrastructure); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot parse input JSON stream",
		})
		return
	}

	infrastructure = importInfrastructure.Infrastructure

	// HTTP 201 Created
	c.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "infrastrcture imported successfully",
	})
}
