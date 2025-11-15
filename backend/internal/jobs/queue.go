package jobs

import (
	"context"
	"fmt"
	"sync"
	
	"github.com/google/uuid"
	"gorm.io/gorm"
	
	"github.com/reserveone/saa-risk-analyzer/internal/domain"
)

type Queue struct {
	db       *gorm.DB
	mu       sync.RWMutex
	jobs     map[uuid.UUID]*JobExecution
	progress map[uuid.UUID]chan int
}

type JobExecution struct {
	Job    *domain.Job
	Func   JobFunc
	Cancel context.CancelFunc
}

func NewQueue(db *gorm.DB) *Queue {
	return &Queue{
		db:       db,
		jobs:     make(map[uuid.UUID]*JobExecution),
		progress: make(map[uuid.UUID]chan int),
	}
}

func (q *Queue) Enqueue(jobType string, fn JobFunc) (*domain.Job, error) {
	job := &domain.Job{
		ID:       uuid.New(),
		Type:     jobType,
		Status:   StatusQueued,
		Progress: 0,
	}
	
	if err := q.db.Create(job).Error; err != nil {
		return nil, err
	}
	
	progressChan := make(chan int, 100)
	
	q.mu.Lock()
	q.jobs[job.ID] = &JobExecution{
		Job:  job,
		Func: fn,
	}
	q.progress[job.ID] = progressChan
	q.mu.Unlock()
	
	go q.executeJob(job.ID)
	
	return job, nil
}

func (q *Queue) executeJob(jobID uuid.UUID) {
	q.mu.RLock()
	exec, ok := q.jobs[jobID]
	progressChan := q.progress[jobID]
	q.mu.RUnlock()
	
	if !ok {
		return
	}
	
	// Update status to running
	exec.Job.Status = StatusRunning
	q.db.Model(exec.Job).Updates(map[string]interface{}{
		"status": StatusRunning,
	})
	
	// Execute function
	result, err := exec.Func(jobID, progressChan)
	
	// Update final status
	if err != nil {
		exec.Job.Status = StatusFailed
		exec.Job.Error = err.Error()
		exec.Job.Progress = 0
	} else {
		exec.Job.Status = StatusSucceeded
		exec.Job.Result = result
		exec.Job.Progress = 100
	}
	
	q.db.Model(exec.Job).Updates(map[string]interface{}{
		"status":   exec.Job.Status,
		"progress": exec.Job.Progress,
		"result":   exec.Job.Result,
		"error":    exec.Job.Error,
	})
	
	close(progressChan)
}

func (q *Queue) GetJob(jobID uuid.UUID) (*domain.Job, error) {
	var job domain.Job
	if err := q.db.First(&job, "id = ?", jobID).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (q *Queue) GetProgress(jobID uuid.UUID) (<-chan int, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	
	ch, ok := q.progress[jobID]
	if !ok {
		return nil, fmt.Errorf("job not found or completed")
	}
	
	return ch, nil
}
