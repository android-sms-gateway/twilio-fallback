package proxy

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	MaxQueueSize = 100
	JobTimeout   = 15 * time.Second
)

type job struct {
	id        string
	operation func(context.Context) error
}

type jobsService struct {
	jobQueue chan job
	workers  int
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc

	logger *zap.Logger
}

func newJobsService(logger *zap.Logger) *jobsService {
	ctx, cancel := context.WithCancel(context.Background())

	return &jobsService{
		jobQueue: make(chan job, MaxQueueSize),
		workers:  runtime.NumCPU() * 2,
		ctx:      ctx,
		cancel:   cancel,

		logger: logger,
	}
}

func (s *jobsService) Start() {
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}
}

func (s *jobsService) Enqueue(id string, operation func(context.Context) error) error {
	select {
	case <-s.ctx.Done():
		return fmt.Errorf("jobs service is closed")
	case s.jobQueue <- job{
		id:        id,
		operation: operation,
	}:
		return nil
	default:
		return fmt.Errorf("job queue full")
	}
}

func (s *jobsService) Close() {
	s.cancel()
	s.wg.Wait() // Wait for all workers to finish
	close(s.jobQueue)
}

func (s *jobsService) worker(id int) {
	defer s.wg.Done()

	for {
		select {
		case job := <-s.jobQueue:
			s.processJob(job, id)
		case <-s.ctx.Done():
			s.logger.Info("Worker stopped", zap.Int("worker", id))
			return
		}
	}
}

func (s *jobsService) processJob(j job, workerID int) {
	s.logger.Info("Processing job", zap.Int("worker", workerID), zap.String("id", j.id))

	ctx, cancel := context.WithTimeout(s.ctx, JobTimeout)
	defer cancel()

	err := j.operation(ctx)
	if err != nil {
		s.logger.Error("Failed to process job",
			zap.Int("worker", workerID),
			zap.String("id", j.id),
			zap.Error(err))
		return
	}

	s.logger.Info("Job processed successfully",
		zap.Int("worker", workerID),
		zap.String("id", j.id))
}
