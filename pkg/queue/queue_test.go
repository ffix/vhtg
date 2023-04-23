package queue_test

//
//import (
//	"fmt"
//	"testing"
//	"time"
//
//	"github.com/golang/mock/gomock"
//
//	"github.com/ffix/vhtg/pkg/queue"
//)
//
//func TestTaskQueue(t *testing.T) {
//	testCases := []struct {
//		name        string
//		task        *queue.Task
//		expectedErr bool
//	}{
//		{
//			name: "Task Success",
//			task: &queue.Task{
//				ID:         "1",
//				Queue:      "A",
//				Payload:    "Task 1 payload",
//				MaxRetries: 3,
//				Timeout:    2 * time.Second,
//				Expiry:     time.Now().Add(6 * time.Second),
//			},
//			expectedErr: false,
//		},
//		{
//			name: "Task Failure",
//			task: &queue.Task{
//				ID:         "2",
//				Queue:      "B",
//				Payload:    "Task 2 payload",
//				MaxRetries: 2,
//				Timeout:    1 * time.Second,
//				Expiry:     time.Now().Add(3 * time.Second),
//			},
//			expectedErr: true,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			mockWorker := NewMockworkerFunc(ctrl)
//			mockLogger := NewMocklogger(ctrl)
//			mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
//
//			taskQueue := queue.NewTaskQueue(mockWorker, mockLogger)
//
//			if tc.expectedErr {
//				mockWorker.EXPECT().Execute(tc.task).Return(fmt.Errorf("task %s failed", tc.task.ID)).AnyTimes()
//			} else {
//				mockWorker.EXPECT().Execute(tc.task).Return(nil).AnyTimes()
//			}
//
//			taskQueue.AddTask(tc.task)
//			taskQueue.WaitAndExit()
//
//		})
//	}
//}
//
//func TestTaskTimeout(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	task := &queue.Task{
//		ID:         "1",
//		Queue:      "A",
//		Payload:    "Task 1 payload",
//		MaxRetries: 1,
//		Timeout:    1 * time.Second,
//		Expiry:     time.Now().Add(5 * time.Second),
//	}
//
//	statusCh := make(chan string, 1)
//
//	mockWorker := NewMockworkerFunc(ctrl)
//	mockWorker.EXPECT().Execute(task).DoAndReturn(func(*queue.Task) error {
//		time.Sleep(2 * time.Second)
//		statusCh <- "timeout"
//		return nil
//	}).Times(1)
//
//	mockLogger := NewMocklogger(ctrl)
//	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
//
//	taskQueue := queue.NewTaskQueue(mockWorker, mockLogger)
//	taskQueue.AddTask(task)
//	taskQueue.WaitAndExit()
//
//	select {
//	case status := <-statusCh:
//		if status != "timeout" {
//			t.Errorf("Expected task to timeout, but it did not")
//		}
//	case <-time.After(5 * time.Second):
//		t.Errorf("Task processing took longer than expected")
//	}
//}
//func TestTaskRetryLimit(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	task := &queue.Task{
//		ID:         "2",
//		Queue:      "B",
//		Payload:    "Task 2 payload",
//		MaxRetries: 2,
//		Timeout:    1 * time.Second,
//		Expiry:     time.Now().Add(5 * time.Second),
//	}
//
//	statusCh := make(chan string, 1)
//
//	mockWorker := NewMockworkerFunc(ctrl)
//	mockWorker.EXPECT().Execute(task).DoAndReturn(func(*queue.Task) error {
//		if task.RetryCount == task.MaxRetries {
//			statusCh <- "retry-limit"
//		}
//		return fmt.Errorf("task %s failed", task.ID)
//	}).Times(task.MaxRetries + 1)
//
//	mockLogger := NewMocklogger(ctrl)
//	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
//
//	taskQueue := queue.NewTaskQueue(mockWorker, mockLogger)
//	taskQueue.AddTask(task)
//	taskQueue.WaitAndExit()
//
//	select {
//	case status := <-statusCh:
//		if status != "retry-limit" {
//			t.Errorf("Expected task to reach retry limit, but it did not")
//		}
//	case <-time.After(5 * time.Second):
//		t.Errorf("Task processing took longer than expected")
//	}
//}
//
//func TestTaskExpiry(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	task := &queue.Task{
//		ID:         "3",
//		Queue:      "C",
//		Payload:    "Task 3 payload",
//		MaxRetries: 3,
//		Timeout:    1 * time.Second,
//		Expiry:     time.Now().Add(3 * time.Second),
//	}
//
//	statusCh := make(chan string, 1)
//
//	mockWorker := NewMockworkerFunc(ctrl)
//	mockWorker.EXPECT().Execute(task).DoAndReturn(func(*queue.Task) error {
//		if time.Now().After(task.Expiry) {
//			statusCh <- "expired"
//		}
//		return fmt.Errorf("task %s failed", task.ID)
//	}).AnyTimes()
//
//	mockLogger := NewMocklogger(ctrl)
//	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
//
//	taskQueue := queue.NewTaskQueue(mockWorker, mockLogger)
//	taskQueue.AddTask(task)
//	taskQueue.WaitAndExit()
//
//	select {
//	case status := <-statusCh:
//		if status != "expired" {
//			t.Errorf("Expected task to be expired, but it did not")
//		}
//	case <-time.After(5 * time.Second):
//		t.Errorf("Task processing took longer than expected")
//	}
//}
