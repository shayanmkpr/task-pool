package main

import (
	"testing"

	"github.com/shayanmkpr/task-pool/internal/models"
	"github.com/shayanmkpr/task-pool/internal/store"
)

// Test basic memory store functionality
func TestMemoryStoreBasic(t *testing.T) {
	store := store.NewMemoryStore()

	task := &models.Task{
		ID:          "basic-test-1",
		Title:       "Basic Test Task",
		Description: "Task for basic store functionality test",
		Status:      models.Pending,
		Duration:    1,
	}

	store.AddTask(task)

	retrievedTask, err := store.GetTask("basic-test-1")
	if err != nil {
		t.Error("Failed to retrieve task:", err)
	}

	if retrievedTask.ID != task.ID {
		t.Error("Retrieved task ID doesn't match")
	}
}

// Test concurrent submissions of many tasks at once (from TODO.md)
func TestMemoryStoreConcurrentSubmissions(t *testing.T) {
	store := store.NewMemoryStore()

	// Add multiple tasks concurrently
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(id int) {
			task := &models.Task{
				ID:          "concurrent-task-" + string(rune(id+'0')),
				Title:       "Concurrent Task " + string(rune(id+'0')),
				Description: "Task for concurrent submission test",
				Status:      models.Pending,
				Duration:    1,
			}
			store.AddTask(task)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all tasks were stored
	tasks := store.ListTasks()
	if len(tasks) != 10 {
		t.Errorf("Expected 10 tasks, got %d", len(tasks))
	}
}

// Test duplicate task IDs (from TODO.md)
func TestMemoryStoreDuplicateIDs(t *testing.T) {
	store := store.NewMemoryStore()

	task1 := &models.Task{
		ID:          "duplicate-test",
		Title:       "First Task",
		Description: "First task with duplicate ID",
		Status:      models.Pending,
		Duration:    1,
	}

	task2 := &models.Task{
		ID:          "duplicate-test", // Same ID
		Title:       "Second Task",
		Description: "Second task with duplicate ID",
		Status:      models.Pending,
		Duration:    1,
	}

	store.AddTask(task1)
	store.AddTask(task2) // This should overwrite the first task

	retrievedTask, err := store.GetTask("duplicate-test")
	if err != nil {
		t.Error("Failed to retrieve task:", err)
	}

	// The second task should overwrite the first one
	if retrievedTask.Title != "Second Task" {
		t.Error("Task was not overwritten properly")
	}
}

// Test empty store
func TestMemoryStoreEmpty(t *testing.T) {
	store := store.NewMemoryStore()

	tasks := store.ListTasks()
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks in empty store, got %d", len(tasks))
	}

	_, err := store.GetTask("non-existent")
	if err == nil {
		t.Error("Expected error when getting non-existent task")
	}
}

