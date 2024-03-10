package workerInfo

import "time"

func IsWorkerStuck(worker *Worker, healthTimeout, progressTimeout time.Duration) bool {
	timeSinceLastHealth := time.Since(worker.LastHealthyTime)
	timeSinceLastProgress := time.Since(worker.LastProgressTime)

	// Check if the worker hasn't responded for a long time
	if timeSinceLastHealth > healthTimeout {
		// Additional check: Has the worker made any progress recently?
		// In place of 20 will set some threshold here
		if timeSinceLastProgress > progressTimeout || worker.Progress < 20 {
			return true // Worker is likely stuck
		}
	}
	return false
}
