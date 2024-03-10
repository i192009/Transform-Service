package workerInfo

import (
	"fmt"
	"strings"
)

// IsWorkerCrashed checks if the worker process has crashed.
func IsWorkerCrashed(worker *Worker) bool {
	if worker.Cmd.ProcessState != nil && worker.Cmd.ProcessState.Exited() {
		return true
	}
	return false
}

// analyzeWorkerCrash fetches and analyzes the logs of the crashed worker
func analyzeWorkerCrash(worker *Worker) {
	// Fetch logs from where the worker was writing them
	// This could be a file, a database, or a logging service
	logs, err := FetchWorkerLogs(worker)
	if err != nil {
		fmt.Println("Error fetching logs for crashed worker:", err)
		return
	}

	// Analyze the logs to determine the potential cause of the crash
	AnalyzeLogs(logs)
}

// fetchWorkerLogs retrieves the logs for the given worker
func FetchWorkerLogs(worker *Worker) ([]string, error) {
	// Implement log retrieval logic here
	// Example: Read from a log file or a central logging database
	return nil, fmt.Errorf("not implemented")
}

func AnalyzeLogs(logs []string) {
	// This is a placeholder for your log analysis logic.
	var foundErrors []string
	for _, log := range logs {
		if isErrorLog(log) {
			foundErrors = append(foundErrors, log)
		}
	}

	if len(foundErrors) > 0 {
		fmt.Println("Found errors in worker logs:")
		for _, errLog := range foundErrors {
			fmt.Println(errLog)
		}
	} else {
		fmt.Println("No errors found in worker logs.")
	}
}

// isErrorLog checks if a log entry indicates an error.
func isErrorLog(log string) bool {
	// Example logic: check if the log contains the word "ERROR"
	return strings.Contains(log, "ERROR")
}
