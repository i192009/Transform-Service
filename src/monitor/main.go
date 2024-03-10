package main

import (
	"fmt"
	"log"
	"math"
	"time"
	"transform2/monitor/workerInfo"
)

// Assuming these are global or part of your system's configuration
var maxQueueLengthPerWorker = 10
var idleTimeThreshold = 5 * time.Minute

const (
	MaxRetries     = 3
	InitialBackoff = 1 * time.Second // Initial delay duration
	BackoffFactor  = 2               // Factor by which to multiply delay for each retry
)

func needMoreWorkers(taskQueueLength, currentWorkerCount int) bool {
	requiredWorkers := (taskQueueLength + maxQueueLengthPerWorker - 1) / maxQueueLengthPerWorker
	return requiredWorkers > currentWorkerCount
}

func tooManyIdleWorkers(workers []*workerInfo.Worker) bool {
	idleWorkers := 0
	for _, worker := range workers {
		if time.Since(worker.LastHealthyTime) > idleTimeThreshold {
			idleWorkers++
		}
	}
	return idleWorkers > len(workers)/2
}

func main() {
	var workers []*workerInfo.Worker
	retryCount := make(map[string]int) // Initialize the retry count map here
	taskQueueLength := getTaskQueueLength()
	// Implement this function
	healthTimeout := 1 * time.Minute
	progressTimeout := 5 * time.Minute

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if needMoreWorkers(taskQueueLength, len(workers)) {
			worker, err := workerInfo.StartWorker("path/to/worker", "NewWorker", "DefaultConfig", 1)
			if err == nil {
				workers = append(workers, worker)
			}
		}

		if tooManyIdleWorkers(workers) {
			// Logic to remove idle workers
			// RemoveWorker is a placeholder for the actual removal logic
			workers = RemoveIdleWorkers(workers)
		}

		for i, worker := range workers {
			go func(w *workerInfo.Worker) {
				if workerInfo.IsWorkerHealthy(w) {
					fmt.Println("Worker is healthy:", w.Name)
					w.LastHealthyTime = time.Now()
				} else {
					fmt.Println("Worker is not healthy:", w.Name)
				}
			}(worker)

			// Fetch logs and analyze the crash
			logs, err := workerInfo.FetchWorkerLogs(worker)
			if err != nil {
				fmt.Println("Error fetching logs:", err)
			} else {
				workerInfo.AnalyzeLogs(logs)
			}

			if workerInfo.IsWorkerCrashed(worker) {
				fmt.Println("Worker has crashed:", worker.Name)
				newWorker, err := retryWorker(worker, retryCount)
				if err != nil {
					fmt.Printf("Failed to restart worker %s: %s\n", worker.Name, err)
					continue
				}
				workers[i] = newWorker
				sendCrashNotification(worker) // Send notification about the crash
				workers = append(workers[:i], workers[i+1:]...)
				i--
				continue
			}

			if workerInfo.IsWorkerStuck(worker, healthTimeout, progressTimeout) {
				fmt.Println("Worker is stuck:", worker.Name)
				fmt.Printf("Attempting to retry worker %s\n", worker.Name)
				newWorker, err := retryWorker(worker, retryCount)
				if err != nil {
					fmt.Printf("Failed to restart worker %s: %s\n", worker.Name, err)
					continue
				}
				workers[i] = newWorker
				// Implement auto-recovery logic here
				workers = append(workers[:i], workers[i+1:]...)
				i--
				continue
			}
		}
	}
}

// Implement these helper functions based on your system's logic
func getTaskQueueLength() int {
	// Return the current length of the task queue
	return 0 // Placeholder
}

func RemoveIdleWorkers(workers []*workerInfo.Worker) []*workerInfo.Worker {
	var activeWorkers []*workerInfo.Worker

	for _, worker := range workers {
		if time.Since(worker.LastHealthyTime) <= idleTimeThreshold {
			activeWorkers = append(activeWorkers, worker)
		} else {
			gracefullyShutdownWorker(worker)
		}
	}
	return activeWorkers
}

func gracefullyShutdownWorker(worker *workerInfo.Worker) {
	// Send a shutdown command or signal to the worker.
	_, _ = worker.Stdin.Write([]byte("shutdown\n"))

	// Alternatively, you could use a signal:
	// _ = worker.Cmd.Process.Signal(os.Interrupt) // or syscall.SIGTERM

	// Wait for the worker to shut down, up to a certain timeout
	timeout := 10 * time.Second
	done := make(chan error, 1)
	go func() {
		done <- worker.Cmd.Wait() // This waits for the worker process to exit
	}()

	select {
	case <-time.After(timeout):
		// Timeout reached, force kill the process
		_ = worker.Cmd.Process.Kill()
	case <-done:
		// Worker shut down gracefully
	}
}

func sendCrashNotification(worker *workerInfo.Worker) {
	// Example: Log the crash event. Replace this with actual notification logic.
	log.Printf("Notification: Worker %s (ID: %s) has crashed.", worker.Name, worker.Id)
}

// retryWorker attempts to restart a worker that has crashed or become stuck.
func retryWorker(worker *workerInfo.Worker, retryCount map[string]int) (*workerInfo.Worker, error) {
	if retryCount[worker.Id] >= MaxRetries {
		return nil, fmt.Errorf("max retries reached for worker %s", worker.Name)
	}

	// Calculate the backoff duration
	backoff := InitialBackoff * time.Duration(math.Pow(BackoffFactor, float64(retryCount[worker.Id])))
	time.Sleep(backoff) // Wait for the backoff duration
	// First, ensure the current worker is shut down.
	gracefullyShutdownWorker(worker)

	// Now, start a new worker.
	newWorker, err := workerInfo.StartWorker(worker.Configuration, worker.Name, worker.Configuration, worker.QueueType)
	if err != nil {
		return nil, err
	}

	// Increment the retry count for this worker.
	retryCount[newWorker.Id] = retryCount[worker.Id] + 1

	return newWorker, nil
}
