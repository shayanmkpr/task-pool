package models

import (
	"time"
)

type Status string

const (
	Pending   Status = "pending"
	Running   Status = "running"
	Completed Status = "completed"
)

type Task struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Duration    time.Duration `json:"duration"` // in seconds
	Status      Status        `json:"status"`
}
