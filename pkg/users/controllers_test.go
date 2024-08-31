package users

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

func TestPostNewUser(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var user User = User{
		ID:        "operator",
		Name:      "operator",
		FullName:  "Mr. Operator",
		TokenHash: "",
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetUsers(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var items = struct {
		Users map[string]User `json:"items"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &items)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, items.Users)
}

func TestGetUserByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/users/operator", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var user = struct {
		User User `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &user)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, user.User)
}

func TestUpdateUserByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var user User = User{
		ID:        "operator",
		Name:      "operator",
		FullName:  "Mrs. Operator",
		TokenHash: "0x33",
		Active:    false,
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("PUT", "/users/operator", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		User User `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, item.User.TokenHash, "0x33")
}

func TestDeleteUserByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("DELETE", "/users/operator", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Key string `json:"key"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "operator", ret.Key)
}

func TestPostDumpRestore(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var items = struct {
		Users map[string]User `json:"items"`
	}{
		Users: map[string]User{
			"operator": {
				ID:     "operator",
				Name:   "operator",
				Active: true,
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
	req, _ := http.NewRequest("POST", "/users/restore", bytes.NewBuffer(jsonValue))

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
	r := core.SetupTestEnv(TestPackage)

	active := true

	req, _ := http.NewRequest("PUT", "/users/operator/active", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		User User `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, !active, item.User.Active)
}

func TestPostUserSSHKeys(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	keys := struct {
		Keys []string `json:"keys"`
	}{
		Keys: []string{
			"ssh-rsa AAAAB3Nzaza42EAAAADAQABAAABgQCr/69RZt3kwGrCkPKt0sP4cQ4z opkey1@station",
			"ssh-rsa AAAAB3Nzaza42EAAAADAQABAAABgQCe/69RZt3kwGrCkPKt0sP4cQ4y opkey2@station",
		},
	}

	jsonValue, _ := json.Marshal(keys)
	req, _ := http.NewRequest("POST", "/users/operator/keys/ssh", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		User User `json:"user"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusAccepted, w.Code)
	assert.Equal(t, keys.Keys, item.User.SSHKeys)
}

func TestGetUserSSHKeys(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/users/operator/keys/ssh", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var keys string = "ssh-rsa AAAAB3Nzaza42EAAAADAQABAAABgQCr/69RZt3kwGrCkPKt0sP4cQ4z opkey1@station" + "\n" + "ssh-rsa AAAAB3Nzaza42EAAAADAQABAAABgQCe/69RZt3kwGrCkPKt0sP4cQ4y opkey2@station"

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, keys, string(w.Body.Bytes()))
}
