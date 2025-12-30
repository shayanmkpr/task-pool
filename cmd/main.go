package main

import (
	"fmt"
	"math/rand"
	"time"

	cfg "github.com/shayanmkpr/task-pool/config"
	"github.com/shayanmkpr/task-pool/internal/models"
	"github.com/shayanmkpr/task-pool/internal/store"
	"github.com/shayanmkpr/task-pool/internal/taskpool"
)

func main() {

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

	for i := 1; i <= 10; i++ {
		task := &models.Task{
			ID:          fmt.Sprintf("task-%d", i),
			Title:       fmt.Sprintf("Task %d", i),
			Description: "Demo task",
			Duration:    time.Duration(rand.Intn(5) + 1), // 1-5 seconds
		}
		pool.AddTask(task)
	}

	memoryStore.WaitForCompletion(500 * time.Millisecond)

	for _, w := range workers {
		w.Stop()
	}
}
