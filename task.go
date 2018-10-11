package main

import (
	"errors"

	"github.com/474420502/grequests"
)

// Task 爬虫个例 以Cookie为个例 或者 以代理IP为个例
type Task struct {
	Config  *TaskConfig
	Session *grequests.Session
}

// NewTask new 一个person 对象
func NewTask(taskFileName string) *Task {

	t := &Task{}
	t.Config = &TaskConfig{}
	t.Config.Load(taskFileName)

	t.Session = grequests.NewSession(nil)

	return t
}

// Execute 更新Session从turl
func (t *Task) Execute() (*grequests.Response, error) {
	t.UpdateSessionFromTURL(t.Config.TURL)
	// spew.Dump(t.Session)
	switch t.Config.TURL.Method {
	case "POST":
		return t.Session.Post(t.Config.TURL.BaseURL)
	case "GET":
		return t.Session.Get(t.Config.TURL.BaseURL)
	case "PATCH":
		return t.Session.Patch(t.Config.TURL.BaseURL)
	case "DELETE":
		return t.Session.Delete(t.Config.TURL.BaseURL)
	case "HEAD":
		return t.Session.Head(t.Config.TURL.BaseURL)
	case "PUT":
		return t.Session.Put(t.Config.TURL.BaseURL)
	case "OPTIONS":
		return t.Session.Options(t.Config.TURL.BaseURL)
	}

	return nil, errors.New("the method is not exists! " + t.Config.TURL.Method)
}

// UpdateSessionFromTURL 更新Session从turl
func (t *Task) UpdateSessionFromTURL(turl *TaskURL) {
	t.Session.RequestOptions.Headers = turl.Headers
	t.Session.RequestOptions.Cookies = turl.Cookies
	t.Session.RequestOptions.Params = turl.Params
	t.Session.RequestOptions.Data = turl.Data

	// log.Println(t.Session.HTTPClient.Transport)
}
