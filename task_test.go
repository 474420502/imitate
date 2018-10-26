package main

import (
	"testing"
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

	t.Run("test execute plan TimeTo interval and at is not TimeTo", func(t *testing.T) {
		for exec, resp := task.ExecuteOnPlan(); exec != nil; {
			t.Error(resp.Content())
		}
	})
}
