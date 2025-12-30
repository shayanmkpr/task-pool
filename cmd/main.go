package main

import (
	"fmt"
	"math/rand"
	"time"

	cfg "github.com/shayanmkpr/task-pool/config"
	"github.com/shayanmkpr/task-pool/internal/logger"
	"github.com/shayanmkpr/task-pool/internal/models"
	"github.com/shayanmkpr/task-pool/internal/store"
	"github.com/shayanmkpr/task-pool/internal/taskpool"
)

func main() {
	lg, err := logger.New("./app.log")
	if err != nil {
		panic(err)
	}
	defer lg.Close()

	lg.Info("Application started")

	config := cfg.NewConfig()
	config.LoadDefaults()

	memoryStore := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(config.PoolSize, memoryStore)

	workers := make([]*taskpool.Worker, config.WorkerCount)
	for i := range config.WorkerCount {
		w := taskpool.NewWorker(i+1, pool)
		w.Start()
		workers[i] = w
	}

	for _, worker := range workers {
		w := worker
		go func() {
			for assigned := range w.Assigned {
				if assigned != nil {
					fmt.Printf("worker %d assigned to %v \n", w.ID, w.ID)
				} else {
					fmt.Printf("worker %d is free %v \n", w.ID, assigned)
				}
			}
		}()
	}

	for i := 1; i <= 10; i++ {
		task := &models.Task{
			ID:          fmt.Sprintf("task-%d", i),
			Title:       fmt.Sprintf("Task %d", i),
			Description: "Demo task",
			Duration:    time.Duration(rand.Intn(5) + 1), // Should be 1-5 seconds
		}
		pool.AddTask(task)
	}

	WaitForCompletion(lg, memoryStore, 500*time.Millisecond)

	for _, w := range workers {
		w.Stop()
	}

	lg.Info("Application finished")
}

func WaitForCompletion(log *logger.Logger, s *store.MemoryStore, waitingTime time.Duration) {
	for {
		allDone := true
		tasks := s.ListTasks()
		for _, t := range tasks {
			if t.Status != models.Completed {
				allDone = false
				break
			}
		}
		if allDone {
			log.Info("All tasks completed")
			return
		}
		time.Sleep(waitingTime)
	}
}
