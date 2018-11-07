package main

import (
	"path/filepath"
)

// TypeMode 多种模式
// 1. 以Proxy为个例
// 2. 以Cookie为个例
type TypeMode int

const (
	_ TypeMode = iota

	// ModeCookie 以Cookie为个例
	ModeCookie = 0
	// ModeProxy  以Proxy为个例
	ModeProxy = 1
)

// Person 拥有两种类型并发任务状态, 一种是以cookie为个例控制, 另一种是以 代理为个例子.
type Person struct {
	Tasks []Task
}

// NewPerson 创建一个新Person
func NewPerson(params ...string) *Person {
	p := Person{}

	for _, tpath := range params {
		p.LoadTasks(tpath)
	}

	return &p
}

// LoadTasks 加载任务
func (person *Person) LoadTasks(tpath string) {
	matches, err := filepath.Glob(tpath)
	if err != nil {
		panic(err)
	}

	for _, match := range matches {
		task := NewTask(match)

		switch task.Config.Setting.Mode {
		case ModeCookie:
			person.Tasks = append(person.Tasks, *task)
		case ModeProxy:
			for _, t := range task.SplitFromProxies() {
				person.Tasks = append(person.Tasks, t)
			}
		}
	}

}

// Execute 执行
func (person *Person) Execute() {
	//TODO: Python的脚本函数, 与动态更新 返回的数据在python的脚本处理
	for _, task := range person.Tasks {
		for _, PResult := range task.ExecuteOnPlan() {
			d := NewPyDict()
			d.UpdateStrStr("content", PResult.Resp.Content())
			d.UpdateStrInt("status", PResult.Resp.GResponse.StatusCode)
			sr := &ScriptResult{NextDo: task.Config.Setting.NextDo, Result: d.PyObject()}
			callScript(sr)
		}
	}
}
