package queue

import (
	"github.com/gin-gonic/gin"
)

// queue CRUD -- functions in controllers.go
func Routes(g *gin.RouterGroup) {
	g.GET("/tasks",
		GetTasks)
	g.POST("/tasks",
		PostNewTask)
	g.GET("/tasks/types",
		ListTypesTasks)
	g.GET("/tasks/:key",
		GetTaskByKey)
	g.PUT("/tasks/:key",
		UpdateTaskByKey)
	g.DELETE("/tasks/:key",
		DeleteTaskByKey)
	g.PUT("/tasks/:key/processed",
		ProcessedToggleByKey)
	g.POST("/restore",
		PostDumpRestore)
}
