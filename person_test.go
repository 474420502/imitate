package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestPersonExecute(t *testing.T) {

	p := NewPerson("task/*_config.py")
	if len(p.Tasks) == 0 {
		t.Error("error load tasks", p)
	}

	os.Remove("/script/save.pyc")
	os.Remove("/script/doothers.pyc")

	c := make(chan bool)
	go func(cr chan bool) {
		time.Sleep(time.Second * 2)
		p.Execute()
		cr <- true
	}(c)

	if <-c {
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
	} else {
		t.Error("c is false?")
	}

}

func TestLoadScrit(t *testing.T) {
	if len(ScriptBook) == 0 {
		t.Error("error load tasks", ScriptBook)
	}
}
