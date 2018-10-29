package main

import (
	"path/filepath"
	"testing"
)

func TestMain(t *testing.T) {
	matches, err := filepath.Glob("task/*_config.py")
	if err != nil {
		panic(err)
	}

	for _, match := range matches {
		task := NewTask(match)
		t.Error("task", task)
	}
}
