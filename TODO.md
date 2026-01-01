# Todo

## Main Tasks:
- [x] Http handlers.
- [x] Http router.
- [x] Graceful Shutdown of the server.
- [ ] Full Configuration.
    - [x] Config factory
    - [ ] All configs added to factory.
- [x] Save logs.
    - [x] remove all log.Printf()s and use slog for them only. fmt.Printf if necessary to have in console.
- [x] Unit tests with a single run all code.
    - [x] Core functionality (workers, pool, memory).
    - [x] Handlers core functionality.
    - [x] Http E2E test. --> used Bruno.
- [x] State Machine?
- [x] Dockerfile and Dockerbuild.
- [x] README with complete instructions for building and running the code.
- [x] (Evaluate solution) Add a coordinator to assign tasks to worker? Over Kill probably.
- [x] For each request, add a timeout. one for all.

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
