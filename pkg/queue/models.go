package queue

import (
	"time"
)

type Task struct {
	// Taks's ID as an unique identifier in the queue.
	ID string `json:"name"`

	// Description of such task in human-readable format.
	Description string `json:"description"`

	// An identifier of a worker for the task to be executed with.
	WorkerName string `json:"worker_name" binding:"required" required:"true"`

	// Processed is an indicator that such task is being worked on.
	Processed bool `json:"processed"`

	// State is a human-readable string describing the actual state of the task.
	State string `json:"state"`

	// Timestamp of the last activity made on this task.
	LastChangeTimestamp time.Time `json:"last_change_timestamp"`
}
