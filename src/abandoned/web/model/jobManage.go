package model

import "sync"

type JobManager struct {
	JobMap map[string]*Job
	RWLock sync.RWMutex
}
