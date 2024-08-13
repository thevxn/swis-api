package links

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

func TestPostNewLink(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var link Link = Link{
		ID:          "sd",
		Name:        "sd",
		Description: "A shortcut for the savla.dev homepage.",
		URL:         "https://savla.dev",
		Active:      true,
	}

	jsonValue, _ := json.Marshal(link)
	req, _ := http.NewRequest("POST", "/links", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	/*var ret = struct {
		Link Link `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &link)*/

	assert.Equal(t, http.StatusCreated, w.Code)
	//assert.Equal(t, link, ret.Link)
}

func TestGetLinks(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/links", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var links = struct {
		Links map[string]Link `json:"items"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &links)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, links.Links)
}

func TestGetLinkByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/links/sd", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var link = struct {
		Link Link `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &link)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, link.Link)
}

func TestUpdateLinkByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var link Link = Link{
		ID:          "sd",
		Name:        "sd",
		Description: "A shortcut for the savla.dev homepage.",
		URL:         "https://www.savla.dev",
		Active:      false,
	}

	jsonValue, _ := json.Marshal(link)
	req, _ := http.NewRequest("PUT", "/links/sd", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Link Link `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, item.Link.Active, false)
}

func TestDeleteLinkByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("DELETE", "/links/sd", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Key string `json:"key"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "sd", ret.Key)
}

func TestPostDumpRestore(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var items = struct {
		Links map[string]Link `json:"items"`
	}{
		Links: map[string]Link{
			"sd": {
				ID:          "sd",
				Name:        "sd",
				Description: "A shortcut for the savla.dev homepage.",
				URL:         "https://www.savla.dev",
				Active:      false,
			},
			/* run #1: this item was 'crippled' on purpose to see how binding would act */
			/* result: it cannot be arsed, all fields are exported to JSON, even unlisted ones... */
			/* --- */
			/* run #2: blank keys SHOULD be ignored at all --- patched in pkg/core/package.go */
			/* result: the struct below is skipped */
			"": {
				ID:          "",
				Name:        "",
				Description: "A blank link",
				URL:         "http://about:blank",
				Active:      false,
			},
		},
	}

	jsonValue, _ := json.Marshal(items)
	req, _ := http.NewRequest("POST", "/links/restore", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Count int `json:"count"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, 1, ret.Count)
}
