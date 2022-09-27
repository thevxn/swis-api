package infra

import (
	//b64 "encoding/base64"
	//"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

var infrastructure = Infrastructure{}

func findHostByHostname(c *gin.Context) (index *int, h *Host) {
	// loop over hosts
	var hosts = infrastructure.Hosts

	for i, h := range hosts {
		if h.HostnameShort == c.Param("hostname") || h.HostnameFQDN == c.Param("hostname") {
			return &i, &h
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "host not found",
	})
	return nil, nil
}

// @Summary Get whole infrastructure
// @Description get all infrastructure details
// @Tags infra
// @Produce  json
// @Success 200 {object} infra.Infrastructure
// @Router /infra [get]
func GetInfrastructure(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":           http.StatusOK,
		"infrastructure": infrastructure,
	})
}

// @Summary Get all hosts
// @Description get hosts list
// @Tags infra
// @Produce  json
// @Success 200 {object} infra.Hosts
// @Router /infra/hosts [get]
func GetHosts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":  http.StatusOK,
		"hosts": infrastructure.Hosts,
	})
}

// @Summary Get all networks
// @Description get networks list
// @Tags infra
// @Produce  json
// @Success 200 {object} infra.Infrastructure.Networks
// @Router /infra/networks [get]
func GetNetworks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"networks": infrastructure.Networks,
	})
}

// @Summary Get host by Hostname
// @Description get host by :hostname param
// @Tags infra
// @Produce  json
// @Success 200 {object} infra.Host
// @Router /infra/hosts/{hostname} [get]
func GetHostByHostname(c *gin.Context) {
	if _, host := findHostByHostname(c); host != nil {
		// host found
		c.IndentedJSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"host": host,
		})
	}
}

// @Summary Upload infrastructure JSON dump
// @Description restore infrastructure data model
// @Tags infra
// @Accept json
// @Produce json
// @Router /infra/restore [post]
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
