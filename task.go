package main

import (
	"errors"

	"github.com/474420502/requests"
)

// Task 爬虫个例 以Cookie为个例 或者 以代理IP为个例
type Task struct {
	Config  *TaskConfig
	Session *requests.Session
}

// NewTask new 一个person 对象
func NewTask(taskFileName string) *Task {

	t := &Task{}

	t.Session = requests.NewSession()

	t.Config = &TaskConfig{}
	t.Config.Load(taskFileName)

	return t
}

// AutoSetSession 从TaskConfig 配置 Session的信息
func (t *Task) AutoSetSession() {

	// TODO:

	// t.Session.SetCookies()
}

// Execute 更新Session从turl
func (t *Task) Execute() (*requests.Response, error) {
	// spew.Dump(t.Session)
	var wf *requests.Workflow
	switch t.Config.Info.Method {
	case "POST":
		wf = t.Session.Post(t.Config.Info.BaseURL)
	case "GET":
		wf = t.Session.Get(t.Config.Info.BaseURL)
	case "PATCH":
		wf = t.Session.Patch(t.Config.Info.BaseURL)
	case "DELETE":
		wf = t.Session.Delete(t.Config.Info.BaseURL)
	case "HEAD":
		wf = t.Session.Head(t.Config.Info.BaseURL)
	case "PUT":
		wf = t.Session.Put(t.Config.Info.BaseURL)
	case "OPTIONS":
		wf = t.Session.Options(t.Config.Info.BaseURL)
	}

	if wf != nil {

	}

	return nil, errors.New("the method is not exists! " + t.Config.Info.Method)
}
