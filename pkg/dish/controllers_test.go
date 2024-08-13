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

// app = array of pointers to pointers to Cache
type appCache []**core.Cache

var TestPackage *core.Package = &core.Package{
	Name: pkgName,
	Cache: []**core.Cache{
		&CacheIncidents,
		&CacheSockets,
	},
	Routes: Routes,
}

// variable used to catch the Incident.ID property
var incidentID string

/*
 *  unit/integration tests
 */

/*
 *  common
 */

func TestPostDumpRestore(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var items = struct {
		Incidents map[string]Incident `json:"incidents"`
		Sockets   map[string]Socket   `json:"sockets"`
	}{
		Incidents: map[string]Incident{
			"outage": {
				Name:     "outage",
				Public:   true,
				SocketID: "",
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
		Sockets: map[string]Socket{
			"socket": {
				ID:   "socket",
				Name: "socket",
			},
		},
	}

	jsonValue, _ := json.Marshal(items)
	req, _ := http.NewRequest("POST", "/dish/restore", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Count []int `json:"counter"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	//t.Logf("%s", jsonValue)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, []int{1, 1}, ret.Count)
}

func TestGetDishRoot(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/dish/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var items = struct {
		Incidents map[string]Incident `json:"incidents"`
		Sockets   map[string]Socket   `json:"sockets"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &items)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, items.Incidents)
	assert.NotEmpty(t, items.Sockets)
}

/*
 *  incidents
 */

func TestPostNewIncident(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	now := time.Now().UnixNano()

	var inc Incident = Incident{
		Name:           "site down",
		StartTimestamp: now,
		Public:         true,
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
	r := core.SetupTestEnv(TestPackage)

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
	r := core.SetupTestEnv(TestPackage)

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
	r := core.SetupTestEnv(TestPackage)

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
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/dish/incidents/non-existent", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Incident Incident `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	// may throw HTTP/500 due to nil CacheSockets while testing incidents' handlers only
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Empty(t, item.Incident)
}

func TestUpdateIncidentByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

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
	r := core.SetupTestEnv(TestPackage)

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

/*
 *  sockets
 */

func TestPostNewSocketByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var socket Socket = Socket{
		ID:   "test_socket",
		Name: "test-socket",
		Host: "host.example.com",
		Port: 80,
		DishTarget: []string{
			"dish_target1",
		},
		Muted:   true,
		Public:  true,
		Healthy: false,
	}

	jsonValue, _ := json.Marshal(socket)
	req, _ := http.NewRequest("POST", "/dish/sockets", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetSockets(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/dish/sockets", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var sockets = struct {
		Sockets map[string]Socket `json:"items"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &sockets)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, sockets.Sockets)
}

func TestGetSocketListPublic(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/dish/sockets/public", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var sockets = struct {
		Sockets map[string]Socket `json:"items"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &sockets)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, sockets.Sockets)
}

/*func TestGetSSEvents(t *testing.T) {}*/

func TestUpdateSocketByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var socket Socket = Socket{
		ID:   "test_socket",
		Name: "test-socket",
		Host: "host.example.com",
		Port: 80,
		DishTarget: []string{
			"dish_target1",
			"dish_target2",
		},
		Muted:   false,
		Public:  true,
		Healthy: false,
	}

	jsonValue, _ := json.Marshal(socket)
	req, _ := http.NewRequest("PUT", "/dish/sockets/test_socket", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Socket Socket `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 2, len(item.Socket.DishTarget))
	assert.Equal(t, false, item.Socket.Muted)
	assert.Equal(t, true, item.Socket.Public)
	assert.Equal(t, false, item.Socket.Healthy)
}

func TestGetSocketListByHost(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/dish/sockets/dish_target2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var sockets = struct {
		Sockets map[string]Socket `json:"items"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &sockets)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, sockets.Sockets)
}

func TestMuteToggleSocketByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	muted := true

	req, _ := http.NewRequest("PUT", "/dish/sockets/test_socket/mute", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Socket Socket `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, !muted, item.Socket.Muted)
}

func TestMaintenanceToggleSocketByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	maintenance := true

	req, _ := http.NewRequest("PUT", "/dish/sockets/test_socket/maintenance", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Socket Socket `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, !maintenance, item.Socket.Maintenance)
}

func TestPublicToggleSocketByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	public := true

	req, _ := http.NewRequest("PUT", "/dish/sockets/test_socket/public", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Socket Socket `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, !public, item.Socket.Public)
}

func TestBatchPostHealthyStatus(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	results := make(map[string]bool)
	results["test_socket"] = true

	socketNames := struct {
		Results map[string]bool `json:"dish_results"`
	}{
		Results: results,
	}

	jsonValue, _ := json.Marshal(socketNames)
	req, _ := http.NewRequest("POST", "/dish/sockets/results", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Count int `json:"count"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, 1, item.Count)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteSocketByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("DELETE", "/dish/sockets/test_socket", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Key string `json:"key"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test_socket", ret.Key)
}
