# Todo

## Main Tasks:
- [ ] Http handlers.
- [ ] Http router.
- [ ] Graceful Shutdown of the server.
- [ ] Full Configuration.
- [ ] Save logs.
- [ ] Unit tests with a single run all code.
- [ ] State Machine?
- [ ] Dockerfile and Dockerbuild.
- [ ] README with complete instructions for building and running the code.

## Edge Cases:
- [ ] The pool is full but the user sends a new task.
- [ ] Pool is empty.
- [ ] Worker count is zero.
- [ ] Task with zero or negative duration.
- [ ] Task with missing or empty title/description.
- [ ] Duplicate task IDs submitted.
- [ ] **Concurrent submissions of many tasks at once.**
- [ ] Worker crashes or panics while processing a task.
- [ ] Graceful shutdown while tasks are still pending.
- [ ] Tasks with extremely long durations.
- [ ] Store fails (simulated error) when adding/updating a task.
- [ ] Tasks blocked due to full buffered channel.
- [ ] Shutdown initiated while new tasks are being added.
- [ ] Memory leak due to unconsumed tasks in the channel.
- [ ] Worker receives a nil task (unexpected input).
- [ ] Long-running tasks delaying other tasks (starvation scenario).
