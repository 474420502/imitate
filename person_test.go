package main

import (
	"testing"
)

func TestPersonLoadTaskConifg(t *testing.T) {
	person := NewPerson("task/task_config")
	if person == nil {
		t.Error("person is nil")
	}

}
