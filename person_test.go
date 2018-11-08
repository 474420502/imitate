package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestPersonExecute(t *testing.T) {
	var err error

	p := NewPerson("task/*_config.py")
	if len(p.Tasks) == 0 {
		t.Error("error load tasks", p)
	}

	os.Remove("/script/save.pyc")
	os.Remove("/script/doothers.pyc")

	time.Sleep(2)
	p.Execute()

	f, err := os.Open("/tmp/test.html")
	if err != nil {
		t.Error(err)
	}

	out, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}

	err = os.Remove("/tmp/test.html")
	if err != nil {
		t.Error(err)
	}

	if strings.LastIndex(string(out), "doothers") == -1 {
		t.Error(string(out), "content error")
	}
}

func TestLoadScrit(t *testing.T) {
	if len(ScriptBook) == 0 {
		t.Error("error load tasks", ScriptBook)
	}
}
