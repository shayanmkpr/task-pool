package task

import "context"

type StoryRepository interface {
	AddTask(task *Task) error
	GetTask(ctx context.Context, id string) (*Task, error)
	ListTasks(ctx context.Context) ([]*Task, error)
	UpdateTask(task *Task)
}
