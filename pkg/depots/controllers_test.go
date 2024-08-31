package depots

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

func TestPostNewDepotItem(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var item DepotItem = DepotItem{
		ID:          "1",
		Description: "An absolutely generic item.",
	}

	jsonValue, _ := json.Marshal(item)
	req, _ := http.NewRequest("POST", "/depots/items", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetAllDepotItems(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/depots/items", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var items = struct {
		Items map[string]DepotItem `json:"items"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &items)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, items.Items)
}

func TestGetDepotItemByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/depots/items/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Item DepotItem `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "1", item.Item.ID)
	assert.NotEmpty(t, item.Item)
}

func TestUpdateDepotItemByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var item DepotItem = DepotItem{
		ID:          "1",
		Description: "An absolutely generic item, edited.",
	}

	jsonValue, _ := json.Marshal(item)
	req, _ := http.NewRequest("PUT", "/depots/items/1", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Item DepotItem `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	// tests
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "1", ret.Item.ID)
	assert.Equal(t, item.Description, ret.Item.Description)
}

func TestDeleteDepotItemByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("DELETE", "/depots/items/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Key string `json:"key"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "1", ret.Key)
}

func TestPostDumpRestore(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var items = struct {
		DepotItem map[string]DepotItem `json:"items"`
	}{
		DepotItem: map[string]DepotItem{
			"1": {
				ID:          "1",
				Description: "An absolutely generic item, edited.",
			},
			/* run #1: this item was 'crippled' on purpose to see how binding would act */
			/* result: it cannot be arsed, all fields are exported to JSON, even unlisted ones... */
			/* --- */
			/* run #2: blank keys SHOULD be ignored at all --- patched in pkg/core/package.go */
			/* result: the struct below is skipped */
			"": {
				ID:          "0",
				Description: "A blank item.",
			},
		},
	}

	jsonValue, _ := json.Marshal(items)
	req, _ := http.NewRequest("POST", "/depots/restore", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Count int `json:"count"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, 1, ret.Count)
}
