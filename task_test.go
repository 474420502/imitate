package main

import (
	"testing"
)

func TestPersonExecute(t *testing.T) {
	person := NewTask("task/task_config")
	if person == nil {
		t.Error("person is nil")
	}

	// ses := grequests.NewSession(nil)
	// resp, err := grequests.Get("https://httpbin.org/gzip", nil)
	resp, err := person.Execute()

	if err != nil || resp.StatusCode != 200 {
		t.Error(resp.Header)
	}
	// f, _ := os.OpenFile("test.json", os.O_CREATE|os.O_RDWR, 777)
	t.Error(resp.DecompressContent)
	t.Error(resp.String())
	// f.WriteString(resp.String())

	//t.Error(person.Config.TURL.Method, person.Config.TURL.BaseURL)
	//t.Error(person.Session.RequestOptions.Params)
}
