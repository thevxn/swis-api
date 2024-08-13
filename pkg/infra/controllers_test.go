package infra

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.savla.dev/swis/v5/pkg/core"

	"github.com/stretchr/testify/assert"
)

// app = array of pointers to pointers to Cache
type appCache []**core.Cache

var TestPackage *core.Package = &core.Package{
	Name: pkgName,
	Cache: []**core.Cache{
		&CacheDomains,
		&CacheHosts,
		&CacheNetworks,
	},
	Routes: Routes,
}

/*
 *  unit/integration tests
 */

/*
 *  common
 */

func TestPostDumpRestore(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var items = struct {
		Domains  map[string]Domain  `json:"Domains"`
		Hosts    map[string]Host    `json:"Hosts"`
		Networks map[string]Network `json:"networks"`
	}{
		Domains: map[string]Domain{
			"test_domain": {
				ID:    "test_domain",
				FQDN:  "example.com",
				Owner: "fantomas",
			},
			"invalid_domain": {
				ID: "",
			},
			"": {
				ID: "",
			},
		},
		Hosts: map[string]Host{
			"test_host": {
				ID:            "test_host",
				HostnameShort: "host",
				HostnameFQDN:  "host.example.com",
			},
			"invalid_host": {
				ID: "",
			},
			"": {
				ID: "",
			},
		},
		Networks: map[string]Network{
			"net_br32": {
				ID:        "net_br32",
				Hash:      "net_br32",
				Name:      "net_br32",
				Interface: "br32",
				CIDRBlock: "/27",
			},
			"invalid_net": {
				ID:   "",
				Hash: "",
			},
			"": {
				ID:   "",
				Hash: "",
			},
		},
	}

	jsonValue, _ := json.Marshal(items)
	req, _ := http.NewRequest("POST", "/infra/restore", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Count []int `json:"counter"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusCreated, w.Code)

	// this is not good, number 2 here means, that the 'invalid' account is imported too,
	// even though it does not suffice given property requirements...
	// TODO: implement simple JSON binding, or some sort of an item verificator
	assert.Equal(t, []int{2, 2, 2}, ret.Count)
}

func TestGetInfrastructure(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/infra", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var items = struct {
		Domains  map[string]Domain  `json:"domains"`
		Hosts    map[string]Host    `json:"hosts"`
		Networks map[string]Network `json:"networks"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &items)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, items.Domains)
	assert.NotEmpty(t, items.Hosts)
	assert.NotEmpty(t, items.Networks)
}

/*
 *  domains
 */

func TestPostNewDomain(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var dom Domain = Domain{
		ID:    "test_domain",
		FQDN:  "example.com",
		Owner: "fantomas",
	}

	jsonValue, _ := json.Marshal(dom)
	req, _ := http.NewRequest("POST", "/infra/domains", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Domain Domain `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	// this account has already been imported/created using the dump restore test
	// --- thus resulting in HTTP/409 Conflict
	//assert.NotEqual(t, "", item.Account.AccountNumber)
	//assert.NotEqual(t, "", item.Account.ID)
	//assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestGetDomains(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/infra/domains", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var items = struct {
		Domains map[string]Domain `json:"items"`
		Count   int               `json:"count"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &items)

	assert.Equal(t, http.StatusOK, w.Code)

	// number 2 down below is wrong, see notes above (cca L87)
	assert.Equal(t, 2, items.Count)
	assert.NotEmpty(t, items.Domains)
}

func TestGetDomainByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/infra/domains/test_domain", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Domain Domain `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, item.Domain)
	assert.NotEqual(t, "", item.Domain.ID)
	assert.NotEqual(t, "", item.Domain.FQDN)
}

func TestUpdateDomainByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var dom Domain = Domain{
		ID:    "test_domain",
		FQDN:  "example.com",
		Owner: "obelix",
	}

	jsonValue, _ := json.Marshal(dom)
	req, _ := http.NewRequest("PUT", "/infra/domains/test_domain", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Domain Domain `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, "", item.Domain.ID)
	assert.NotEqual(t, "", item.Domain.FQDN)
	assert.Equal(t, "obelix", item.Domain.Owner)
}

func TestDeleteDomainByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("DELETE", "/infra/domains/test_domain", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Key string `json:"key"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test_domain", ret.Key)
}

/*
 *  hosts
 */

func TestPostNewHost(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var host Host = Host{
		ID:            "test_host",
		HostnameShort: "host",
		HostnameFQDN:  "host.example.com",
	}

	jsonValue, _ := json.Marshal(host)
	req, _ := http.NewRequest("POST", "/infra/hosts", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	//assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestGetHosts(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/infra/hosts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var items = struct {
		Hosts map[string]Host `json:"items"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &items)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, items.Hosts)
}

func TestGetHostByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/infra/hosts/test_host", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Host Host `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, item.Host)
}

func TestUpdateHostByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var host Host = Host{
		ID:            "test_host",
		HostnameShort: "host",
		HostnameFQDN:  "host.prod.example.com",
	}

	jsonValue, _ := json.Marshal(host)
	req, _ := http.NewRequest("PUT", "/infra/hosts/test_host", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Host Host `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, host.HostnameFQDN, ret.Host.HostnameFQDN)
}

func TestDeleteHostByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("DELETE", "/infra/hosts/test_host", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Key string `json:"key"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test_host", ret.Key)
}

/*
 *  networks
 */

func TestPostNewNetwork(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var net Network = Network{
		ID:        "net_br32",
		Hash:      "net_br32",
		Name:      "net_br32",
		Interface: "br32",
		CIDRBlock: "/27",
	}

	jsonValue, _ := json.Marshal(net)
	req, _ := http.NewRequest("POST", "/infra/networks", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	//assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestGetNetworks(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/infra/networks", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var nets = struct {
		Networks map[string]Network `json:"items"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &nets)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, nets.Networks)
}

func TestGetNetworkByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/infra/networks/net_br32", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var net = struct {
		Network Network `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &net)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, net.Network)
}

func TestUpdateNetworkByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var net Network = Network{
		ID:        "net_br32",
		Hash:      "net_br32",
		Name:      "net_br32",
		Interface: "br32",
		CIDRBlock: "/26",
	}

	jsonValue, _ := json.Marshal(net)
	req, _ := http.NewRequest("PUT", "/infra/networks/net_br32", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Network Network `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, "", item.Network.Hash)
	assert.Equal(t, "/26", item.Network.CIDRBlock)
}

func TestDeleteNetworkByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("DELETE", "/infra/networks/net_br32", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Key string `json:"key"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "net_br32", ret.Key)
}
