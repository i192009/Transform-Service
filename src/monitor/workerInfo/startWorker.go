package workerInfo

import (
	"github.com/google/uuid"
	"io"
	"os/exec"
	"time"
)

type Worker struct {
	Cmd              *exec.Cmd      // Holds the *exec.Cmd object
	Stdin            io.WriteCloser // Standard input of the worker process
	Stdout           io.ReadCloser  // Standard output of the worker process
	LastHealthyTime  time.Time      // Time of the last successful health check
	Id               string         // Unique Identifier
	Name             string         // Name of the Worker
	Status           bool           // Whether it's busy or Free
	Configuration    string         // ID of the configuration it's using
	QueueType        int64          // Type of the Queue it is listening to (VIP or Normal)
	LastProgressTime time.Time      // Time of the last known progress
	Progress         float64        // Progress percentage or similar metric
}

// StartWorker starts the worker process and sets up stdin and stdout pipes.
func StartWorker(cmdPath, name, config string, queueType int64) (*Worker, error) {
	cmd := exec.Command(cmdPath)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Worker{
		Cmd:              cmd,
		Stdin:            stdin,
		Stdout:           stdout,
		LastHealthyTime:  time.Now(),
		Id:               generateUniqueID(), // Implement this function to generate unique IDs
		Name:             name,
		Status:           true,       // Initially busy
		LastProgressTime: time.Now(), // Initialize the last progress time
		Progress:         0.0,        // Initialize progress to 0%Configuration:   config,
		QueueType:        queueType,
	}, nil
}

// Utility function to generate unique ID for a worker
func generateUniqueID() string {
	return uuid.New().String()
}
