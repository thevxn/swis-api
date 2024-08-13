package backups

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

func TestPostBackeduServicepByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var bcp Backup = Backup{
		ID:          "swapi",
		ServiceName: "swapi",
		Description: "A very swapi service.",
		Active:      true,
	}

	jsonValue, _ := json.Marshal(bcp)
	req, _ := http.NewRequest("POST", "/backups", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetBackups(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/backups", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var backups = struct {
		Backups map[string]Backup `json:"items"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &backups)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, backups.Backups)
}

func TestGetBackedupStatusByServiceKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/backups/swapi", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var bcp = struct {
		Backup Backup `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &bcp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, bcp.Backup)
}

func TestUpdateBackupStatusByServiceKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var bcp Backup = Backup{
		ID:          "swapi",
		ServiceName: "swapi",
		Description: "A very swapi service.",
		Active:      false,
	}

	jsonValue, _ := json.Marshal(bcp)
	req, _ := http.NewRequest("PUT", "/backups/swapi", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Backup Backup `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, false, item.Backup.Active)
}

func TestDeleteBackupByServiceKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("DELETE", "/backups/swapi", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Key string `json:"key"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "swapi", ret.Key)
}

func TestPostDumpRestore(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var items = struct {
		Backups map[string]Backup `json:"items"`
	}{
		Backups: map[string]Backup{
			"swapi": {
				ID:          "swapi",
				ServiceName: "swapi",
				Description: "A very swapi service.",
				Active:      true,
			},
			/* run #1: this item was 'crippled' on purpose to see how binding would act */
			/* result: it cannot be arsed, all fields are exported to JSON, even unlisted ones... */
			/* --- */
			/* run #2: blank keys SHOULD be ignored at all --- patched in pkg/core/package.go */
			/* result: the struct below is skipped */
			"": {
				ID:          "",
				ServiceName: "",
				Description: "A blank service Name.",
				Active:      false,
			},
		},
	}

	jsonValue, _ := json.Marshal(items)
	req, _ := http.NewRequest("POST", "/backups/restore", bytes.NewBuffer(jsonValue))

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

func TestActiveToggleBackupByServiceKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	active := true

	req, _ := http.NewRequest("PUT", "/backups/swapi/active", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Backup Backup `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, !active, item.Backup.Active)
}
