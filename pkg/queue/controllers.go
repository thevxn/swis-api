package queue

import (
	"net/http"
	"strconv"
	"time"

	"go.savla.dev/swis/v5/pkg/core"

	"github.com/gin-gonic/gin"
)

var (
	CacheTasks *core.Cache
	pkgName    string = "queue"
)

var Package *core.Package = &core.Package{
	Name: pkgName,
	Cache: []**core.Cache{
		&CacheTasks,
	},
	Routes: Routes,
}

// GetLinks returns JSON serialized list of tasks and their properties.
// @Summary Get all tasks
// @Description get tasks complete list
// @Tags queue
// @Produce json
// @Success 200 {object} queue.Task
// @Router /queue/tasks [get]
func GetTasks(ctx *gin.Context) {
	core.PrintAllRootItems(ctx, CacheTasks, pkgName)
	return
}

// GetTaskByKey returns task's properties, given sent key exists in database.
// @Summary Get task by timestamp key
// @Description get task by its key param
// @Tags queue
// @Produce json
// @Success 200 {object} queue.Task
// @Router /queue/tasks/{key} [get]
func GetTaskByKey(ctx *gin.Context) {
	core.PrintItemByParam[Task](ctx, CacheTasks, pkgName, Task{})
	return
}

// PostNewTask produces a new task
//
// @Summary      Add new task to queue
// @Description  add new task to queue
// @Tags         queue
// @Accept       json
// @Produce      json
// @Success      201  {object}  queue.Task
// @Failure      400  {object}  queue.Task
// @Failure      409  {object}  queue.Task
// @Failure      500  {object}  queue.Task
// @Router       /queue/tasks [post]
func PostNewTask(ctx *gin.Context) {
	var id string

	// loop until new task.ID is generated and usable
	for {
		id = strconv.FormatInt(time.Now().UnixNano(), 10)

		_, found := CacheTasks.Get(id)
		if !found {
			break
		}
	}

	var newTask Task

	if err := ctx.BindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
			"message": "cannot parse input JSON stream",
		})
		return
	}

	newTask.ID = id
	newTask.LastChangeTimestamp = time.Now()

	if saved := CacheTasks.Set(id, newTask); !saved {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "task couldn't be saved to database",
		})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "new task created",
		"task":    newTask,
		"id":      id,
	})
	return
}

// @Summary Upload tasks dump backup -- restore all tasks
// @Description update tasks JSON dump
// @Tags queue
// @Accept json
// @Produce json
// @Router /queue/restore [post]
func PostDumpRestore(ctx *gin.Context) {
	core.BatchRestoreItems[Task](ctx, CacheTasks, pkgName)
	return
}

// @Summary Update task by its Key
// @Description update task by its Key
// @Tags queue
// @Produce json
// @Param request body queue.Task.ID true "query params"
// @Success 200 {object} queue.Task
// @Router /queue/tasks/{key} [put]
func UpdateTaskByKey(ctx *gin.Context) {
	core.UpdateItemByParam[Task](ctx, CacheTasks, pkgName, Task{})
	return
}

// @Summary Delete task by its Key
// @Description delete task by its Key
// @Tags queue
// @Produce json
// @Param  id  path  string  true  "task Key"
// @Success 200 {object} queue.Task
// @Router /queue/tasks/{key} [delete]
func DeleteTaskByKey(ctx *gin.Context) {
	core.DeleteItemByParam(ctx, CacheTasks, pkgName)
	return
}

// @Summary Toggle processed boolean by task's key
// @Description toggle processed boolean by task's key
// @Tags queue
// @Produce json
// @Param  id  path  string  true  "key"
// @Success 200 {object} queue.Task
// @Router /queue/tasks/{key}/processed [put]
func ProcessedToggleByKey(ctx *gin.Context) {
	var key string = ctx.Param("key")
	var task Task

	rawTask, ok := CacheTasks.Get(key)
	if !ok {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "task not found",
		})
		return
	}

	task, ok = rawTask.(Task)
	if !ok {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "cannot assert data type, database internal error",
		})
		return
	}

	// inverse the Processed field value
	task.Processed = !task.Processed
	task.LastChangeTimestamp = time.Now()

	if saved := CacheTasks.Set(key, task); !saved {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "task couldn't be saved to database",
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "task processed toggle pressed!",
		"task":    task,
	})
	return
}
