package main

import (
	"errors"

	"github.com/474420502/requests"
)

// Task 爬虫个例 以Cookie为个例 或者 以代理IP为个例
type Task struct {
	Config  *TaskConfig
	Session *requests.Session
	Proxies []string
}

// NewTask new 一个person 对象
func NewTask(taskFileName string) *Task {

	t := &Task{}

	t.Session = requests.NewSession()

	t.Config = &TaskConfig{}
	t.Config.Load(taskFileName)

	t.AutoSetSession()

	return t
}

// AutoSetSession 从TaskConfig 配置 Session的信息
func (t *Task) AutoSetSession() {

	// TODO:
	if t.Config != nil {
		t.Session.Query = t.Config.Info.Query
		t.Session.Header = t.Config.Info.Header
		t.Proxies = t.Config.Setting.Proxies
		t.Session.SetConfig(requests.ConfigRequestTimeout, 120)
	}

	//t.Session.SetCookies()
}

// SplitFromProxies 从这个拆分的Task 是没办法自动reload配置
func (t *Task) SplitFromProxies() []Task {
	var result []Task

	for _, proxy := range t.Config.Setting.Proxies {
		tempTask := Task{}
		tempTask.Session = requests.NewSession()

		tempTask.Session.Query = t.Config.Info.Query
		tempTask.Session.Header = t.Config.Info.Header
		tempTask.Proxies = append(tempTask.Proxies, proxy)

		result = append(result, tempTask)
	}

	return result
}

// ExecuteOnPlan 按时执行
func (t *Task) ExecuteOnPlan() []PlanResult {
	var result []PlanResult
	for _, exec := range t.Config.Setting.Plan.ExecuteQueue {
		if exec.TimeTo() >= 0 {
			resp, err := t.Execute()
			if err == nil {
				result = append(result, PlanResult{Exec: exec, Resp: resp})
			}
		}
	}
	return result
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
		if t.Config.Info.Cookies != nil {
			for _, c := range t.Config.Info.Cookies {
				wf.AddCookie(c)
			}
		}

		if t.Proxies != nil {
			t.Session.SetConfig(requests.ConfigProxy, t.Proxies[0])
			t.Proxies = append(t.Proxies[1:], t.Proxies[0])
		}

		return wf.Execute()
	}

	return nil, errors.New("the method is not exists! " + t.Config.Info.Method)
}
