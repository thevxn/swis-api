package infra

import (
	"net/http"

	"go.savla.dev/swis/v5/pkg/core"

	"github.com/gin-gonic/gin"
)

var (
	CacheHosts    *core.Cache
	CacheNetworks *core.Cache
	CacheDomains  *core.Cache
	pkgName       string = "infra"
)

// @Summary Get whole infrastructure
// @Description get all infrastructure details
// @Tags infra
// @Produce  json
// @Success 200 {object} infra.Infrastructure
// @Router /infra [get]
func GetInfrastructure(ctx *gin.Context) {
	domains, _ := CacheDomains.GetAll()
	hosts, _ := CacheHosts.GetAll()
	networks, _ := CacheNetworks.GetAll()

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"message":  "ok, dumping infrastructure",
		"domains":  domains,
		"hosts":    hosts,
		"networks": networks,
	})
}

/*

  DOMAINS CRUD

*/

// @Summary Get all domains
// @Description get domain list
// @Tags infra
// @Produce  json
// @Success 200 {object} infra.Domain
// @Router /infra/domains [get]
func GetDomains(ctx *gin.Context) {
	core.PrintAllRootItems(ctx, CacheDomains, pkgName)
	return
}

// @Summary Get domain by Key
// @Description get domain by :key param
// @Tags infra
// @Produce  json
// @Success 200 {object} infra.Domain
// @Router /infra/domains/{key} [get]
func GetDomainByKey(ctx *gin.Context) {
	core.PrintItemByParam(ctx, CacheDomains, pkgName, Domain{})
	return
}

// @Summary Add new domain
// @Description add new domain
// @Tags infra
// @Produce json
// @Param request body infra.Domain true "query params"
// @Success 200 {object} infra.Domain
// @Router /infra/domains/{key} [post]
func PostNewDomainByKey(ctx *gin.Context) {
	core.AddNewItemByParam(ctx, CacheDomains, pkgName, Domain{})
	return
}

// @Summary Update domain by its Key
// @Description update domain by its Key
// @Tags infra
// @Produce json
// @Param request body infra.Domain.ID true "query params"
// @Success 200 {object} infra.Domain
// @Router /infra/domains/{key} [put]
func UpdateDomainByKey(ctx *gin.Context) {
	core.UpdateItemByParam(ctx, CacheDomains, pkgName, Domain{})
	return
}

// @Summary Delete domain by its Key
// @Description delete domain by its Key
// @Tags infra
// @Produce json
// @Param id path string true "domain ID/Key"
// @Success 200 {object} infra.Domain.ID
// @Router /infra/domains/{key} [delete]
func DeleteDomainByKey(ctx *gin.Context) {
	core.DeleteItemByParam(ctx, CacheDomains, pkgName)
	return
}

/*

  HOSTS CRUD

*/

// @Summary Get all hosts
// @Description get hosts list
// @Tags infra
// @Produce  json
// @Success 200 {object} infra.Host
// @Router /infra/hosts [get]
func GetHosts(ctx *gin.Context) {
	core.PrintAllRootItems(ctx, CacheHosts, pkgName)
	return
}

// @Summary Get host by Key
// @Description get host by :key param
// @Tags infra
// @Produce  json
// @Success 200 {object} infra.Host
// @Router /infra/hosts/{key} [get]
func GetHostByKey(ctx *gin.Context) {
	core.PrintItemByParam(ctx, CacheHosts, pkgName, Host{})
	return
}

// @Summary Upload current host facts
// @Description update host's facts
// @Tags infra
// @Produce json
// @Param request body infra.Host true "query params"
// @Success 200 {object} infra.Host
// @Router /infra/hosts/{key}/facts [post]
func PostHostFactsByKey(ctx *gin.Context) {
	return
}

// @Summary Add new host
// @Description add new host
// @Tags infra
// @Produce json
// @Param request body infra.Host true "query params"
// @Success 200 {object} infra.Host
// @Router /infra/hosts/{key} [post]
func PostNewHostByKey(ctx *gin.Context) {
	core.AddNewItemByParam(ctx, CacheHosts, pkgName, Host{})
	return
}

// @Summary Update host by its Key
// @Description update host by its Key
// @Tags infra
// @Produce json
// @Param request body infra.Host.ID true "query params"
// @Success 200 {object} infra.Host
// @Router /infra/hosts/{key} [put]
func UpdateHostByKey(ctx *gin.Context) {
	core.UpdateItemByParam(ctx, CacheHosts, pkgName, Host{})
	return
}

// @Summary Delete host by its Key
// @Description delete host by its Key
// @Tags infra
// @Produce json
// @Param id path string true "host ID/Key"
// @Success 200 {object} infra.Host.ID
// @Router /infra/hosts/{key} [delete]
func DeleteHostByKey(ctx *gin.Context) {
	core.DeleteItemByParam(ctx, CacheHosts, pkgName)
	return
}

/*

  NETWORKS CRUD

*/

// @Summary Get all networks
// @Description get networks list
// @Tags infra
// @Produce  json
// @Success 200 {object} infra.Infrastructure.Networks
// @Router /infra/networks [get]
func GetNetworks(ctx *gin.Context) {
	core.PrintAllRootItems(ctx, CacheNetworks, pkgName)
	return
}

// @Summary Get network by Key
// @Description get network by :key param
// @Tags infra
// @Produce  json
// @Success 200 {object} infra.Network
// @Router /infra/networks/{key} [get]
func GetNetworkByKey(ctx *gin.Context) {
	core.PrintItemByParam(ctx, CacheNetworks, pkgName, Network{})
	return
}

// @Summary Add new network
// @Description add new network
// @Tags infra
// @Produce json
// @Param request body infra.Netowrk true "query params"
// @Success 200 {object} infra.Network
// @Router /infra/networks/{key} [post]
func PostNewNetworkByKey(ctx *gin.Context) {
	core.AddNewItemByParam(ctx, CacheNetworks, pkgName, Network{})
	return
}

// @Summary Update network by its Key
// @Description update network by its Key
// @Tags infra
// @Produce json
// @Param request body infra.Network.ID true "query params"
// @Success 200 {object} infra.Network
// @Router /infra/networks/{key} [put]
func UpdateNetworkByKey(ctx *gin.Context) {
	core.UpdateItemByParam(ctx, CacheNetworks, pkgName, Network{})
	return
}

// @Summary Delete network by its Key
// @Description delete network by its Key
// @Tags infra
// @Produce json
// @Param id path string true "network ID/Key"
// @Success 200 {object} infra.Network.ID
// @Router /infra/networks/{key} [delete]
func DeleteNetworkByKey(ctx *gin.Context) {
	core.DeleteItemByParam(ctx, CacheNetworks, pkgName)
	return
}

/*

  RESTORATION

*/

// @Summary Upload infrastructure JSON dump
// @Description restore infrastructure data model
// @Tags infra
// @Accept json
// @Produce json
// @Router /infra/restore [post]
// PostDumpRestore
func PostDumpRestore(ctx *gin.Context) {
	var counter []int = []int{0, 0, 0}

	var importInfra = struct {
		Domains  map[string]Domain  `json:"domains"`
		Hosts    map[string]Host    `json:"hosts"`
		Networks map[string]Network `json:"networks"`
	}{}

	if err := ctx.BindJSON(&importInfra); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
			"message": "cannot parse input JSON stream",
		})
		return
	}

	for key, item := range importInfra.Domains {
		CacheDomains.Set(key, item)
		counter[0]++
	}

	for key, item := range importInfra.Hosts {
		CacheHosts.Set(key, item)
		counter[1]++
	}

	for key, item := range importInfra.Networks {
		CacheNetworks.Set(key, item)
		counter[2]++
	}

	// HTTP 201 Created
	ctx.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"counter": counter,
		"message": "infrastrcture imported successfully",
	})
}
