package taskpool

import (
	"log"
	"time"

	"github.com/shayanmkpr/task-pool/internal/models"
)

type Worker struct {
	ID       int
	TaskPool *TaskPool
	Quit     chan struct{}
	Assigned chan *models.Task
}

func NewWorker(id int, pool *TaskPool) *Worker {
	return &Worker{
		ID:       id,
		TaskPool: pool,
		Quit:     make(chan struct{}),
		Assigned: make(chan *models.Task, 1),
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
				log.Printf("Worker %d shutting down", w.ID)
				return
			}
		}
	}()
}

func (w *Worker) process(task *models.Task) {
	// log.Printf("Worker %d started task %s", w.ID, task.ID)
	task.Status = models.Running
	w.TaskPool.Store.UpdateTask(task)
	w.Assigned <- task
	time.Sleep(task.Duration * time.Second)
	task.Status = models.Completed
	w.TaskPool.Store.UpdateTask(task)
	w.Assigned <- nil
	// log.Printf("Worker %d completed task %s", w.ID, task.ID)
}

func (w *Worker) Stop() {
	close(w.Quit)
}
