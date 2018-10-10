package main

import (
	"testing"
)

func TestPersonExecute(t *testing.T) {
	person := NewTask("task/task_config")
	if person == nil {
		t.Error("person is nil")
	}
	resp, err := person.Execute()
	if err != nil || resp.StatusCode != 200 {
		t.Error(resp.String())
	}
	t.Error(resp.String())
}
