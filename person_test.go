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

	for _, match := range matches {
		task := NewTask(match)
		t.Error("task", task)
	}
}
