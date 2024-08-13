package alvax

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.savla.dev/swis/v5/pkg/core"

	"github.com/stretchr/testify/assert"
)

var TestPackage *core.Package = &core.Package{
	Name: pkgName,
	Cache: []**core.Cache{
		&Cache,
	},
	Routes: Routes,
}

/*
 *  unit/integration tests
 */

func TestPostNewConfig(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var cfg ConfigRoot = ConfigRoot{
		ID:  "bot",
		Key: "bot",
	}

	jsonValue, _ := json.Marshal(cfg)
	req, _ := http.NewRequest("POST", "/alvax/", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetConfigs(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/alvax/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var cfgs = struct {
		Configs map[string]ConfigRoot `json:"items"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &cfgs)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, cfgs.Configs)
}

func TestGetConfigByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/alvax/bot", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var cfg = struct {
		Config ConfigRoot `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &cfg)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, cfg.Config)
}

func TestUpdateConfigByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var cfg ConfigRoot = ConfigRoot{
		Key: "bot",
	}

	jsonValue, _ := json.Marshal(cfg)
	req, _ := http.NewRequest("PUT", "/alvax/bot", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Config ConfigRoot `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, item.Config)
}

func TestDeleteConfigByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("DELETE", "/alvax/bot", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Key string `json:"key"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "bot", ret.Key)
}

func TestPostDumpRestore(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var items = struct {
		Configs map[string]ConfigRoot `json:"items"`
	}{
		Configs: map[string]ConfigRoot{
			"bot": {
				Key: "bot",
			},
			/* run #1: this item was 'crippled' on purpose to see how binding would act */
			/* result: it cannot be arsed, all fields are exported to JSON, even unlisted ones... */
			/* --- */
			/* run #2: blank keys SHOULD be ignored at all --- patched in pkg/core/package.go */
			/* result: the struct below is skipped */
			"": {
				Key: "",
			},
		},
	}

	jsonValue, _ := json.Marshal(items)
	req, _ := http.NewRequest("POST", "/alvax/restore", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Count int `json:"count"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	//t.Logf("%s", jsonValue)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, 1, ret.Count)
}
