package roles

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

func TestPostNewRole(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var role Role = Role{
		ID:          "operators",
		Name:        "operators",
		Description: "A very role for operators",
		Admin:       true,
		Active:      true,
	}

	jsonValue, _ := json.Marshal(role)
	req, _ := http.NewRequest("POST", "/roles", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetRoles(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/roles/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var roles = struct {
		Roles map[string]Role `json:"items"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &roles)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, roles.Roles)
}

func TestGetRoleByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/roles/operators", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var role = struct {
		Role Role `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &role)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, role.Role)
}

func TestUpdateRoleByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var role Role = Role{
		ID:        "operators",
		Name:        "operators",
		Description: "A very role for operators",
		Admin:       false,
		Active:      false,
	}

	jsonValue, _ := json.Marshal(role)
	req, _ := http.NewRequest("PUT", "/roles/operators", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Role Role `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, item.Role.Admin, false)
	assert.Equal(t, item.Role.Active, false)
}

func TestDeleteRoleByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("DELETE", "/roles/operators", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Key string `json:"key"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "operators", ret.Key)
}

func TestPostDumpRestore(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var items = struct {
		Roles map[string]Role `json:"items"`
	}{
		Roles: map[string]Role{
			"operators": {
				ID:        "operators",
				Name:        "operators",
				Description: "A very role for operators",
				Admin:       true,
				Active:      true,
			},
			/* run #1: this item was 'crippled' on purpose to see how binding would act */
			/* result: it cannot be arsed, all fields are exported to JSON, even unlisted ones... */
			/* --- */
			/* run #2: blank keys SHOULD be ignored at all --- patched in pkg/core/package.go */
			/* result: the struct below is skipped */
			"": {
				Name:        "",
				Description: "blank role",
				Admin:       true,
				Active:      true,
			},
		},
	}

	jsonValue, _ := json.Marshal(items)
	req, _ := http.NewRequest("POST", "/roles/restore", bytes.NewBuffer(jsonValue))

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
