package workerInfo

import (
	"bufio"
	"encoding/json"
	"fmt"
)

// HealthMetrics struct to hold CPU and Memory Usage health metrics
type HealthMetrics struct {
	CPUUsage    float64 `json:"cpuUsage"`
	MemoryUsage float64 `json:"memoryUsage"`
}

func IsWorkerHealthy(worker *Worker) bool {
	// First, perform the simple alive check
	if !simpleHealthCheck(worker) {
		return false
	}

	// Then, perform the metrics-based health check
	return metricsBasedHealthCheck(worker)
}

// simpleHealthCheck sends a basic health check command and expects a specific response
func simpleHealthCheck(worker *Worker) bool {
	healthCheckCmd := "check_health\n"

	_, err := worker.Stdin.Write([]byte(healthCheckCmd))
	if err != nil {
		return false
	}

	scanner := bufio.NewScanner(worker.Stdout)
	if scanner.Scan() {
		response := scanner.Text()
		return response == "healthy"
	}

	return false
}

// metricsBasedHealthCheck retrieves and evaluates health metrics from the worker
func metricsBasedHealthCheck(worker *Worker) bool {
	_, err := worker.Stdin.Write([]byte("get_metrics\n"))
	if err != nil {
		return false
	}

	var metrics HealthMetrics
	err = json.NewDecoder(worker.Stdout).Decode(&metrics)
	if err != nil {
		fmt.Println("Error decoding metrics:", err)
		return false
	}

	return evaluateHealthMetrics(metrics)
}

// EvaluateHealthMetrics checks if the provided metrics are within acceptable thresholds
func evaluateHealthMetrics(metrics HealthMetrics) bool {
	const MaxCPUUsage = 80.0    // Maximum acceptable CPU usage percentage
	const MaxMemoryUsage = 70.0 // Maximum acceptable memory usage percentage

	// Check if CPU and memory usage are within their respective thresholds
	if metrics.CPUUsage > MaxCPUUsage {
		fmt.Printf("High CPU usage detected: %.2f%%\n", metrics.CPUUsage)
		return false
	}

	if metrics.MemoryUsage > MaxMemoryUsage {
		fmt.Printf("High Memory usage detected: %.2f%%\n", metrics.MemoryUsage)
		return false
	}

	return true
}
