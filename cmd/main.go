package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/shayanmkpr/task-pool/internal/models"
	"github.com/shayanmkpr/task-pool/internal/store"
	"github.com/shayanmkpr/task-pool/internal/taskpool"
)

func main() {
	// should move to config
	poolSize := 10
	workerCount := 3
	memoryStore := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(poolSize, memoryStore)
	workers := make([]*taskpool.Worker, workerCount)
	for i := range workerCount {
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
