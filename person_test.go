package main

import (
	"path/filepath"
	"testing"
)

func TestPerson(t *testing.T) {
	matches, err := filepath.Glob("task/*_config.py")
	if err != nil {
		panic(err)
	}

	var tasks []*Task
	for _, match := range matches {
		task := NewTask(match)
		t.Error("task", task)
		tasks = append(tasks, task)
	}

}
