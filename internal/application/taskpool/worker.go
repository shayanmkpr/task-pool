package taskpool

import (
	"fmt"
	"time"

	taskEntity "github.com/shayanmkpr/task-pool/internal/domain/task"
)

type Worker struct {
	ID       int
	TaskPool *TaskPool
	Quit     chan struct{}
	Assigned chan *taskEntity.Task
}

func NewWorker(id int, pool *TaskPool) *Worker {
	return &Worker{
		ID:       id,
		TaskPool: pool,
		Quit:     make(chan struct{}),
		Assigned: make(chan *taskEntity.Task, 1),
	}
}

func (w *Worker) Start() {
	go func() {
		defer close(w.Assigned)
		for {
			select {
			case task := <-w.TaskPool.Tasks: // reading from a buffered channel. This handles the Queue logic.
				w.process(task)
			case <-w.Quit:
				// close(w.Assigned) // close the assigned channel.
				fmt.Printf("Worker %d shutting down\n", w.ID) //fix
				return
			}
		}
	}()
}

func (w *Worker) process(task *taskEntity.Task) {
	defer func() { // not sure
		if r := recover(); r != nil {
			task.Status = taskEntity.StatusFailed
			w.TaskPool.Store.UpdateTask(task)
			fmt.Printf("Worker %d: task %s failed with panic: %v\n", w.ID, task.ID, r)
		}
	}()
	task.Status = taskEntity.StatusRunning
	w.TaskPool.Store.UpdateTask(task)
	w.Assigned <- task
	time.Sleep(time.Duration(task.Duration) * time.Second)

	task.Status = taskEntity.StatusCompleted
	w.TaskPool.Store.UpdateTask(task)
	w.Assigned <- nil
	fmt.Printf("Worker %d completed task %s\n", w.ID, task.ID) //fix
}

func (w *Worker) Stop() {
	close(w.Quit)
}
