package main

import (
	"testing"
	"time"
)

func TestTaskExecute(t *testing.T) {
	task := NewTask("task/task_config")
	if task == nil {
		t.Error("task is nil")
	}

	t.Run("test base execute", func(t *testing.T) {
		resp, err := task.Execute()
		if err != nil {
			t.Error(resp.Content())
		}
		if len(resp.Content()) <= 500 {
			t.Error(resp.Content())
		}
	})

}

func TestTaskExecuteOnPlan(t *testing.T) {
	task := NewTask("task/task_config")
	if task == nil {
		t.Error("task is nil")
	}

	t.Run("test plan interval", func(t *testing.T) {
		time.Sleep(time.Second * 5)
		if len(task.ExecuteOnPlan()) < 1 {
			t.Error("ExecuteOnPlan is error, Maybe TimeTo ...")
		}
	})
}
