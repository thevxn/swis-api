package queue

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
		&CacheTasks,
	},
	Routes: Routes,
}

/*
 *  unit/integration tests
 */

// helper variable to track task's ID
var taskID string

func TestPostNewTask(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var task Task = Task{
		Description: "A testing task.",
		WorkerName:  "test",
		Processed:   false,
		State:       "new",
	}

	jsonValue, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/queue/tasks", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Task Task `json:"task"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	if ret.Task.ID != "" {
		taskID = ret.Task.ID
	}

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotEmpty(t, ret.Task.ID)
	assert.NotEmpty(t, ret.Task.LastChangeTimestamp)
}

func TestGetTasks(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/queue/tasks", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var items = struct {
		Tasks map[string]Task `json:"items"`
		Count int             `json:"count"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &items)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 1, items.Count)
	assert.NotEmpty(t, items.Tasks)
}

func TestGetTaskByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/queue/tasks/"+taskID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var task = struct {
		Task Task `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &task)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, task.Task)
}

func TestUpdateTaskByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var task Task = Task{
		Description: "A testing task (updated)",
		WorkerName:  "test",
		State:       "updated",
	}

	jsonValue, _ := json.Marshal(task)
	req, _ := http.NewRequest("PUT", "/queue/tasks/"+taskID, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Task Task `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, task.State, item.Task.State)
	assert.Equal(t, task.Description, item.Task.Description)
}

func TestDeleteTaskByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("DELETE", "/queue/tasks/"+taskID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Key string `json:"key"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, taskID, ret.Key)
}

func TestPostDumpRestore(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var items = struct {
		Tasks map[string]Task `json:"items"`
	}{
		Tasks: map[string]Task{
			"123456": {
				Description: "A testing task.",
				WorkerName:  "test",
				State:       "new",
			},
			/* run #1: this item was 'crippled' on purpose to see how binding would act */
			/* result: it cannot be arsed, all fields are exported to JSON, even unlisted ones... */
			/* --- */
			/* run #2: blank keys SHOULD be ignored at all --- patched in pkg/core/package.go */
			/* result: the struct below is skipped */
			"": {
				Description: "A blank task",
			},
		},
	}

	jsonValue, _ := json.Marshal(items)
	req, _ := http.NewRequest("POST", "/queue/restore", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Count int `json:"count"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, 1, ret.Count)
}
