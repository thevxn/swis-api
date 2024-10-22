package business

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.vxn.dev/swis/v5/pkg/core"

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

func TestPostNewBusiness(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var biz Business = Business{
		ID:        "vxn-dev",
		NameLabel: "vxn.dev",
	}

	jsonValue, _ := json.Marshal(biz)
	req, _ := http.NewRequest("POST", "/business", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetBusinessEntities(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/business", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var biz = struct {
		Entities map[string]Business `json:"items"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &biz)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, biz.Entities)
}

func TestGetBusinessByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/business/vxn-dev", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var biz = struct {
		Entity Business `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &biz)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, biz.Entity)
	assert.NotEmpty(t, biz.Entity.ID)
}

func TestUpdateBusinessByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var biz Business = Business{
		ID:        "vxn-dev",
		NameLabel: "vxn.dev ltd.",
	}

	jsonValue, _ := json.Marshal(biz)
	req, _ := http.NewRequest("PUT", "/business/vxn-dev", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Entity Business `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, biz.NameLabel, item.Entity.NameLabel)
	assert.NotEmpty(t, item.Entity)
	assert.NotEmpty(t, item.Entity.ID)
}

func TestDeleteBusinessByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("DELETE", "/business/vxn-dev", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Key string `json:"key"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "vxn-dev", ret.Key)
}

func TestPostDumpRestore(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var items = struct {
		Entities map[string]Business `json:"items"`
	}{
		Entities: map[string]Business{
			"vxn-dev": {
				ID:        "vxn-dev",
				NameLabel: "vxn.dev",
			},
			/* run #1: this item was 'crippled' on purpose to see how binding would act */
			/* result: it cannot be arsed, all fields are exported to JSON, even unlisted ones... */
			/* --- */
			/* run #2: blank keys SHOULD be ignored at all --- patched in pkg/core/package.go */
			/* result: the struct below is skipped */
			"": {
				ID:        "",
				NameLabel: "",
			},
		},
	}

	jsonValue, _ := json.Marshal(items)
	req, _ := http.NewRequest("POST", "/business/restore", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Count []int `json:"count"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, []int{1}, ret.Count)
}
