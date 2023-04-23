package queue

import (
	"sync"
	"time"

	"github.com/ffix/vhtg/pkg/events"
)

const (
	retryTimeout  = 1 * time.Second
	backoffFactor = 2
	channelSize   = 10000
)

type Task struct {
	Payload    events.Event
	Expiry     time.Time
	retryCount int
}

type TaskQueue struct {
	queues     map[int]chan *Task
	workerFunc func(*Task) error
	wg         sync.WaitGroup
	logger     logger
}

func NewTaskQueue(workerFunc func(*Task) error, queueIDs []int, logger logger) *TaskQueue {
	queues := make(map[int]chan *Task)
	for _, queueID := range queueIDs {
		queues[queueID] = make(chan *Task, channelSize)
	}

	tq := &TaskQueue{
		queues:     queues,
		workerFunc: workerFunc,
		logger:     logger,
	}
	for queue := range queues {
		go tq.startWorker(queue)
	}
	return tq
}

func (tq *TaskQueue) AddTask(payload events.Event, expiry time.Time) {
	task := Task{
		Payload:    payload,
		Expiry:     expiry,
		retryCount: 0,
	}

	for _, queue := range tq.queues {
		tq.wg.Add(1)
		queue <- &task
	}
}

func (tq *TaskQueue) startWorker(queue int) {
	for task := range tq.queues[queue] {
		tq.processTask(task)
	}
}

func (tq *TaskQueue) processTask(task *Task) {
	defer tq.wg.Done()

	var err error
	for {
		if time.Now().After(task.Expiry) {
			tq.logger.Warn("Task discarded due to exceeding the expiry time")
			return
		}

		err = tq.workerFunc(task)
		if err == nil {
			return
		}

		task.retryCount++

		// Calculate the time to wait before the next retry
		retryDuration := retryTimeout * time.Duration(backoffFactor<<(task.retryCount-1))

		if retryDuration > time.Until(task.Expiry) {
			tq.logger.Warn("Task discarded due to the next attempt will exceed the expiry time")
			return
		}

		// Wait for the minimum duration between the retry and expiry
		<-time.After(retryDuration)
	}
}

func (tq *TaskQueue) WaitAndExit() {
	tq.wg.Wait()
}
