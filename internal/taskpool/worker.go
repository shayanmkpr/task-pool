package taskpool

import (
	"log"
	"time"

	"github.com/shayanmkpr/task-pool/internal/models"
)

type Worker struct {
	ID       int
	TaskPool *TaskPool
	Quit     chan bool
}

func NewWorker(id int, pool *TaskPool) *Worker {
	return &Worker{
		ID:       id,
		TaskPool: pool,
		Quit:     make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case task := <-w.TaskPool.Tasks:
				w.process(task)
			case <-w.Quit:
				log.Printf("Worker %d shutting down", w.ID)
				return
			}
		}
	}()
}

func (w *Worker) process(task *models.Task) {
	log.Printf("Worker %d started task %s", w.ID, task.ID)
	task.Status = models.Running
	w.TaskPool.Store.UpdateTask(task)

	time.Sleep(task.Duration * time.Second)

	task.Status = models.Completed
	w.TaskPool.Store.UpdateTask(task)
	log.Printf("Worker %d completed task %s", w.ID, task.ID)
}

func (w *Worker) Stop() {
	w.Quit <- true
}
