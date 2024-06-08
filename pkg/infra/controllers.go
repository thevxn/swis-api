package infra

import (
	"net/http"
	"os"

	"go.savla.dev/swis/v5/pkg/core"

	"github.com/gin-gonic/gin"
)

var (
	CacheHosts    *core.Cache
	CacheNetworks *core.Cache
	CacheDomains  *core.Cache
	pkgName       string = "infra"
)

var Package *core.Package = &core.Package{
	Name: pkgName,
	Cache: []**core.Cache{
		&CacheDomains,
		&CacheHosts,
		&CacheNetworks,
	},
	Routes: Routes,
}

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

// @Summary Post domain deployment by key
// @Description post domain deployment by key
// @Tags infra
// @Produce json
// @Param request body []infra.DNSRecord true "query params"
// @Success 200 {object} infra.Domain
// @Router /infra/domains/{key}/deployment [post]
func PostDomainDeploymentByKey(ctx *gin.Context) {
	key := ctx.Param("key")

	rawDomain, ok := CacheDomains.Get(key)
	if !ok {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "domain not found by key",
			"package": pkgName,
			"key":     key,
		})
		return
	}

	var domain Domain

	domain, ok = rawDomain.(Domain)
	if !ok {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "cannot assert Domain data type",
			"package": pkgName,
			"key":     key,
		})
		return
	}

	email := os.Getenv("CF_API_EMAIL")
	token := os.Getenv("CF_API_TOKEN")

	if email == "" || token == "" {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Cloudflare API key and e-mail not provided as ENV variables",
			"package": pkgName,
		})
		return
	}

	var records []DNSRecord

	if err := ctx.BindJSON(&records); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
			"key":     key,
			"message": "cannot bind input JSON stream",
			"package": pkgName,
		})
		return
	}

	// range records to call the Cloudflare API once per a record
	for id, record := range records {
		//if err := callCfAPI([]string{email, token}, domain.CfZoneID, record); err != nil {
		if err := callCfAPI(domain.CfZoneID, record); err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
				"code":      http.StatusInternalServerError,
				"error":     err.Error(),
				"key":       key,
				"record_id": id,
				"message":   "error occurred while calling Cloudflare API",
				"package":   pkgName,
			})
			return
		}
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"key":     key,
		"message": "domain records successfully deployed",
		"package": pkgName,
	})
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

// @Summary Upload current host configuration
// @Description update host's configuration
// @Tags infra
// @Produce json
// @Param request body infra.Configuration true "host's configuration"
// @Success 200 {object} infra.Host
// @Router /infra/hosts/{key}/config [post]
func PostHostConfigByKey(ctx *gin.Context) {
	key := ctx.Param("key")
	rawHost, ok := CacheHosts.Get(key)
	config := Configuration{}

	// look up the host
	if !ok {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "such host not found",
			"package": pkgName,
			"key":     key,
		})
		return
	}

	// assert type Host to the fetched raw data
	host, ok := rawHost.(Host)
	if !ok {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"key":     key,
			"message": "cannot assert data type, database internal error",
			"package": pkgName,
		})
		return
	}

	// load the payload into configuration struct, must bind
	if err := ctx.BindJSON(&config); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
			"message": "cannot parse input JSON stream",
			"package": pkgName,
			"key":     key,
		})
		return
	}

	// overwright the configuration
	host.Configuration = config

	if saved := CacheHosts.Set(key, host); !saved {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"key":     key,
			"message": "item couldn't be saved to database",
			"package": pkgName,
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"item":    host,
		"key":     key,
		"message": "host's configuration updated",
		"packege": pkgName,
	})
	return
}

// @Summary Add/update a VM install configuration
// @Description add/update a VM install configuration
// @Tags infra
// @Produce json
// @Param request body infra.VMInstallConfig true "host's VMIC"
// @Success 200 {object} infra.Host
// @Router /infra/hosts/{key}/vmic [post]
func PostHostVMICByKey(ctx *gin.Context) {
	key := ctx.Param("key")
	rawHost, ok := CacheHosts.Get(key)
	config := VMInstallConfig{}

	// look up the host
	if !ok {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "such host not found",
			"package": pkgName,
			"key":     key,
		})
		return
	}

	// assert type Host to the fetched raw data
	host, ok := rawHost.(Host)
	if !ok {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"key":     key,
			"message": "cannot assert data type, database internal error",
			"package": pkgName,
		})
		return
	}

	// load the payload into VMIC struct, must bind
	if err := ctx.BindJSON(&config); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
			"message": "cannot parse input JSON stream",
			"package": pkgName,
			"key":     key,
		})
		return
	}

	// loop over children, add/overwrite their install configs
	for idx, name := range host.Children {
		if config.LocalName == name {
			configs := host.ChildrenConfigs
			if len(configs) == 0 {
				configs = append(configs, config)
			} else {
				configs[idx] = config
			}
			host.ChildrenConfigs = configs
			break
		}
	}

	// save updated Host struct to backend
	if saved := CacheHosts.Set(key, host); !saved {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"key":     key,
			"message": "item couldn't be saved to database",
			"package": pkgName,
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"item":    host,
		"key":     key,
		"message": "host's VMIC updated",
		"vm":      config.Name,
		"packege": pkgName,
	})
	return
}

// @Summary Delete a VM install configuration
// @Description delete a VM install configuration
// @Tags infra
// @Produce json
// @Success 200 {object} infra.Host
// @Router /infra/hosts/{key}/vmic/{vm} [delete]
func DeleteHostVMICByKeyAndVM(ctx *gin.Context) {
	key := ctx.Param("key")
	vm := ctx.Param("vm")
	rawHost, ok := CacheHosts.Get(key)

	// look up the host
	if !ok {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "such host not found",
			"package": pkgName,
			"key":     key,
		})
		return
	}

	// assert type Host to the fetched raw data
	host, ok := rawHost.(Host)
	if !ok {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"key":     key,
			"message": "cannot assert data type, database internal error",
			"package": pkgName,
		})
		return
	}

	// loop over children, search for key with requested VM name
	for idx, name := range host.Children {
		if vm == name {
			configs := host.ChildrenConfigs
			configs[idx] = VMInstallConfig{}
			host.ChildrenConfigs = configs
		}
	}

	// save the updated Host struct to backend
	if saved := CacheHosts.Set(key, host); !saved {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"key":     key,
			"message": "item couldn't be saved to database",
			"package": pkgName,
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"item":    host,
		"key":     key,
		"message": "host's VMIC deleted",
		"vm":      vm,
		"packege": pkgName,
	})
	return
}

// @Summary Upload current host facts
// @Description update host's facts
// @Tags infra
// @Produce json
// @Param request body infra.Facts true "host's facts"
// @Success 200 {object} infra.Host
// @Router /infra/hosts/{key}/facts [post]
func PostHostFactsByKey(ctx *gin.Context) {
	key := ctx.Param("key")
	rawHost, ok := CacheHosts.Get(key)
	facts := Facts{}

	// look up the host
	if !ok {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "such host not found",
			"package": pkgName,
			"key":     key,
		})
		return
	}

	// assert type Host to the fetched raw data
	host, ok := rawHost.(Host)
	if !ok {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"key":     key,
			"message": "cannot assert data type, database internal error",
			"package": pkgName,
		})
		return
	}

	// load the payload into facts struct, must bind
	if err := ctx.BindJSON(&facts); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
			"message": "cannot parse input JSON stream",
			"package": pkgName,
			"key":     key,
		})
		return
	}

	// overwright the facts
	host.Facts = facts

	if saved := CacheHosts.Set(key, host); !saved {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"key":     key,
			"message": "item couldn't be saved to database",
			"package": pkgName,
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"item":    host,
		"key":     key,
		"message": "host's facts updated",
		"packege": pkgName,
	})
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
// @Param request body infra.Network true "query params"
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
		if key == "" {
			continue
		}

		CacheDomains.Set(key, item)
		counter[0]++
	}

	for key, item := range importInfra.Hosts {
		if key == "" {
			continue
		}

		CacheHosts.Set(key, item)
		counter[1]++
	}

	for key, item := range importInfra.Networks {
		if key == "" {
			continue
		}

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
