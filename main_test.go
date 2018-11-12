package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	p := NewPerson("task/*_config.py")
	if len(p.Tasks) == 0 {
		t.Error("error load tasks", p)
	}
}
