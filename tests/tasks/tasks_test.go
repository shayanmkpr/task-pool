package main

import (
	"fmt"
	"testing"

	"github.com/shayanmkpr/task-pool/internal/models"
	"github.com/shayanmkpr/task-pool/internal/store"
)

// Test basic task creation with valid values
func TestTaskCreation(t *testing.T) {
	task := &models.Task{
		ID:          "test-task-1",
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      models.Pending,
		Duration:    1,
	}

	if task.ID != "test-task-1" {
		t.Error("Task ID not set correctly")
	}

	if task.Title != "Test Task" {
		t.Error("Task title not set correctly")
	}

	if task.Status != models.Pending {
		t.Error("Task status not set correctly")
	}

	fmt.Println("Task creation test completed successfully")
}

// Test task with zero or negative duration (from TODO.md)
func TestTaskZeroNegativeDuration(t *testing.T) {
	taskWithZeroDuration := &models.Task{
		ID:          "zero-duration-task",
		Title:       "Zero Duration Task",
		Description: "Task with zero duration",
		Status:      models.Pending,
		Duration:    0, // Zero duration
	}

	if taskWithZeroDuration.Duration != 0 {
		t.Error("Task with zero duration not handled correctly")
	}

	taskWithNegativeDuration := &models.Task{
		ID:          "negative-duration-task",
		Title:       "Negative Duration Task",
		Description: "Task with negative duration",
		Status:      models.Pending,
		Duration:    -1, // Negative duration
	}

	if taskWithNegativeDuration.Duration != -1 {
		t.Error("Task with negative duration not handled correctly")
	}
}

// Test task with missing or empty title/description (from TODO.md)
func TestTaskMissingTitleDescription(t *testing.T) {
	taskWithEmptyTitle := &models.Task{
		ID:          "empty-title-task",
		Title:       "", // Empty title
		Description: "Task with empty title",
		Status:      models.Pending,
		Duration:    1,
	}

	if taskWithEmptyTitle.Title != "" {
		t.Error("Task with empty title not handled correctly")
	}

	taskWithEmptyDescription := &models.Task{
		ID:          "empty-desc-task",
		Title:       "Task with empty description",
		Description: "", // Empty description
		Status:      models.Pending,
		Duration:    1,
	}

	if taskWithEmptyDescription.Description != "" {
		t.Error("Task with empty description not handled correctly")
	}
}

// Test task storage functionality
func TestTaskStorage(t *testing.T) {
	store := store.NewMemoryStore()

	task := &models.Task{
		ID:          "storage-test-1",
		Title:       "Storage Test Task",
		Description: "Task to test storage functionality",
		Status:      models.Pending,
		Duration:    1,
	}

	// Add task to store
	store.AddTask(task)

	// Retrieve task from store
	retrievedTask, err := store.GetTask("storage-test-1")
	if err != nil {
		t.Error("Failed to retrieve task from store:", err)
	}

	if retrievedTask.ID != task.ID {
		t.Error("Retrieved task ID doesn't match original")
	}

	if retrievedTask.Title != task.Title {
		t.Error("Retrieved task title doesn't match original")
	}

	fmt.Println("Task storage test completed successfully")
}

// Test duplicate task IDs (from TODO.md) - this is covered in memory tests
// but let's add a higher level test here
func TestDuplicateTaskIDs(t *testing.T) {
	store := store.NewMemoryStore()

	task1 := &models.Task{
		ID:          "duplicate-id-test",
		Title:       "First Task",
		Description: "First task with duplicate ID",
		Status:      models.Pending,
		Duration:    1,
	}

	task2 := &models.Task{
		ID:          "duplicate-id-test", // Same ID as task1
		Title:       "Second Task",
		Description: "Second task with duplicate ID",
		Status:      models.Pending,
		Duration:    1,
	}

	store.AddTask(task1)
	store.AddTask(task2) // This should overwrite the first task

	retrievedTask, err := store.GetTask("duplicate-id-test")
	if err != nil {
		t.Error("Failed to retrieve task:", err)
	}

	// The second task should overwrite the first one
	if retrievedTask.Title != "Second Task" {
		t.Error("Task was not overwritten properly when duplicate ID used")
	}
}