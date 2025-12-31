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

	workerManager := taskpool.NewWorkerManager(config.WorkerCount, memoryStore)
	workerManager.InitiateWorkers(pool)
	workers := workerManager.Workers

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

	workerManager.WaitForCompletion(lg, 500*time.Millisecond)

	for _, w := range workers {
		w.Stop()
	}

	lg.Info("Application finished")
}
