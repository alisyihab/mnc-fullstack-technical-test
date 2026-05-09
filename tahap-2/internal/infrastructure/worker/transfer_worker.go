package worker

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// TransferJob holds the data needed to complete an async transfer credit.
type TransferJob struct {
	SenderID      uuid.UUID
	ReceiverID    uuid.UUID
	Amount        float64
	TransactionID uuid.UUID
}

// JobRecord is a snapshot of a completed (or failed) job kept for the dashboard.
type JobRecord struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	SenderID      uuid.UUID `json:"sender_id"`
	ReceiverID    uuid.UUID `json:"receiver_id"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	ProcessedAt   time.Time `json:"processed_at"`
	Error         string    `json:"error,omitempty"`
}

// CreditFn is the function the worker calls to credit the receiver and update
// the transaction status.  Injected at construction so the worker has no import
// cycle with the usecase / repository layers.
type CreditFn func(ctx context.Context, receiverID uuid.UUID, amount float64, transactionID uuid.UUID, status string) error

// TransferWorker is an in-memory, channel-based worker pool for async transfer
// credit operations.
type TransferWorker struct {
	queue       chan TransferJob
	creditFn    CreditFn
	workerCount int

	// atomic stats
	jobsPending   int64
	jobsProcessed int64
	jobsFailed    int64

	// recent job history (max 100, circular)
	mu      sync.Mutex
	history []JobRecord
}

const maxHistory = 100

// NewTransferWorker creates a new TransferWorker.
// workerCount: number of goroutines in the pool.
// queueSize:   capacity of the buffered channel.
func NewTransferWorker(workerCount, queueSize int, creditFn CreditFn) *TransferWorker {
	return &TransferWorker{
		queue:       make(chan TransferJob, queueSize),
		creditFn:    creditFn,
		workerCount: workerCount,
		history:     make([]JobRecord, 0, maxHistory),
	}
}

// Start launches the worker goroutines.  It blocks until ctx is cancelled.
func (w *TransferWorker) Start(ctx context.Context) {
	var wg sync.WaitGroup
	for i := 0; i < w.workerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			w.runWorker(ctx, id)
		}(i)
	}
	wg.Wait()
}

func (w *TransferWorker) runWorker(ctx context.Context, workerID int) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-w.queue:
			if !ok {
				return
			}
			atomic.AddInt64(&w.jobsPending, -1)
			w.processJob(ctx, job, workerID)
		}
	}
}

func (w *TransferWorker) processJob(ctx context.Context, job TransferJob, workerID int) {
	record := JobRecord{
		TransactionID: job.TransactionID,
		SenderID:      job.SenderID,
		ReceiverID:    job.ReceiverID,
		Amount:        job.Amount,
		ProcessedAt:   time.Now(),
	}

	err := w.creditFn(ctx, job.ReceiverID, job.Amount, job.TransactionID, "SUCCESS")
	if err != nil {
		atomic.AddInt64(&w.jobsFailed, 1)
		record.Status = "FAILED"
		record.Error = err.Error()
		log.Printf("[worker-%d] failed to process transfer job %s: %v", workerID, job.TransactionID, err)
	} else {
		atomic.AddInt64(&w.jobsProcessed, 1)
		record.Status = "SUCCESS"
		log.Printf("[worker-%d] processed transfer job %s", workerID, job.TransactionID)
	}

	w.appendHistory(record)
}

func (w *TransferWorker) appendHistory(r JobRecord) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.history) >= maxHistory {
		// drop oldest entry
		w.history = w.history[1:]
	}
	w.history = append(w.history, r)
}

// Enqueue adds a job to the queue.  Returns false if the queue is full.
func (w *TransferWorker) Enqueue(job TransferJob) bool {
	select {
	case w.queue <- job:
		atomic.AddInt64(&w.jobsPending, 1)
		return true
	default:
		log.Printf("transfer worker queue full, dropping job %s", job.TransactionID)
		return false
	}
}

// Stats returns a snapshot of queue statistics.
func (w *TransferWorker) Stats() map[string]any {
	return map[string]any{
		"jobs_pending":   atomic.LoadInt64(&w.jobsPending),
		"jobs_processed": atomic.LoadInt64(&w.jobsProcessed),
		"jobs_failed":    atomic.LoadInt64(&w.jobsFailed),
		"worker_count":   w.workerCount,
	}
}

// RecentJobs returns a copy of the recent job history.
func (w *TransferWorker) RecentJobs() []JobRecord {
	w.mu.Lock()
	defer w.mu.Unlock()
	result := make([]JobRecord, len(w.history))
	copy(result, w.history)
	return result
}
