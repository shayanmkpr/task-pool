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
	workerManager.MonitorWorkers(lg)
	workerManager.ForceStopWorkers() // should change later to the time limit instead of defer
	workerManager.WaitForCompletion(lg, 500*time.Millisecond)

	// here will be replaces with http server

	for i := 1; i <= 10; i++ {
		task := &models.Task{
			ID:          fmt.Sprintf("task-%d", i),
			Title:       fmt.Sprintf("Task %d", i),
			Description: "Demo task",
			Duration:    time.Duration(rand.Intn(5) + 1), // Should be 1-5 seconds
		}
		pool.AddTask(task)
	}

	lg.Info("Application finished")
}
