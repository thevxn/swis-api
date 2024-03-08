package dish

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.savla.dev/swis/v5/pkg/core"

	"github.com/stretchr/testify/assert"
)

var TestPackageInc *core.Package = &core.Package{
	Name:   pkgName,
	Cache:  &CacheIncidents,
	Routes: Routes,
}

var TestPackageSoc *core.Package = &core.Package{
	Name:   pkgName,
	Cache:  &CacheSockets,
	Routes: Routes,
}

// variable used to catch the Incident.ID property
var incidentID string

/*
 *  unit/integration tests
 */

/*
 *  incidents
 */

func TestPostNewIncident(t *testing.T) {
	r := core.SetupTestEnv(TestPackageInc)

	now := time.Now().UnixNano()

	var inc Incident = Incident{
		ID:             "123123",
		Name:           "site down",
		StartTimestamp: now,
		Public:         true,
		//SocketID:  "https_savla_dev",
	}

	jsonValue, _ := json.Marshal(inc)
	req, _ := http.NewRequest("POST", "/dish/incidents", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		ID string `json:"id"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	if item.ID != "" {
		incidentID = item.ID
	}

	assert.NotEqual(t, "", item.ID)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetIncidentList(t *testing.T) {
	r := core.SetupTestEnv(TestPackageInc)

	req, _ := http.NewRequest("GET", "/dish/incidents", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var items = struct {
		Incidents map[string]Incident `json:"items"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &items)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, items.Incidents)
}

func TestGetGlobalIncidentList(t *testing.T) {
	r := core.SetupTestEnv(TestPackageInc)

	req, _ := http.NewRequest("GET", "/dish/incidents/global", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var items = struct {
		//Incidents map[string]Incident `json:"items"`
		Incidents []Incident `json:"items"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &items)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, items.Incidents)
}

func TestGetPublicIncidentList(t *testing.T) {
	r := core.SetupTestEnv(TestPackageInc)

	req, _ := http.NewRequest("GET", "/dish/incidents/public", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var items = struct {
		Incidents []Incident `json:"items"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &items)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, items.Incidents)
}

func TestGetIncidentListBySocketID(t *testing.T) {
	r := core.SetupTestEnv(TestPackageInc)

	req, _ := http.NewRequest("GET", "/dish/incidents/non-existent", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Incident Incident `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	// HTTP/500 due to nil CacheSockets while testing incidents' handlers
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Empty(t, item.Incident)
}

func TestUpdateIncidentByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackageInc)

	now := time.Now().UnixNano()
	then := time.Now().Add(time.Hour * 1).UnixNano()

	var inc Incident = Incident{
		Name:           "site down",
		StartTimestamp: now,
		EndTimestamp:   then,
		Public:         false,
		//SocketID:  "https_savla_dev",
	}

	jsonValue, _ := json.Marshal(inc)
	req, _ := http.NewRequest("PUT", "/dish/incidents/"+incidentID, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Incident Incident `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, false, item.Incident.Public)
	assert.Equal(t, now, item.Incident.StartTimestamp)
	assert.Equal(t, then, item.Incident.EndTimestamp)
}

func TestDeleteIncidentByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackageInc)

	req, _ := http.NewRequest("DELETE", "/dish/incidents/"+incidentID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Key string `json:"key"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, incidentID, ret.Key)
}

func TestPostDumpRestore(t *testing.T) {
	r := core.SetupTestEnv(TestPackageInc)

	var items = struct {
		Incidents map[string]Incident `json:"items"`
	}{
		Incidents: map[string]Incident{
			"operator": {
				Name: "operator",
			},
			/* run #1: this item was 'crippled' on purpose to see how binding would act */
			/* result: it cannot be arsed, all fields are exported to JSON, even unlisted ones... */
			/* --- */
			/* run #2: blank keys SHOULD be ignored at all costs --- patched in pkg/core/package.go */
			/* result: the project struct below is skipped */
			"": {
				Name: "",
			},
		},
	}

	jsonValue, _ := json.Marshal(items)
	req, _ := http.NewRequest("POST", "/dish/restore", bytes.NewBuffer(jsonValue))

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

func TestActiveToggleUserByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackageInc)

	active := true

	req, _ := http.NewRequest("PUT", "/users/operator/active", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Socket Socket `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, !active, item.Socket.Muted)
}
