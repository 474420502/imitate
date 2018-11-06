package main

import (
	"testing"
)

func TestPersonExecute(t *testing.T) {
	p := NewPerson("task/*_config.py")
	if len(p.Tasks) == 0 {
		t.Error("error load tasks", p)
	}
	p.Execute()
	t.Error(1)
}

func TestLoadScrit(t *testing.T) {
	if len(ScriptBook) == 0 {
		t.Error("error load tasks", ScriptBook)
	}
}
