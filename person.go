package main

import (
	"path/filepath"

	"github.com/474420502/requests"
	py "github.com/sbinet/go-python"
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
func NewPerson(tpath string) *Person {
	p := Person{}

	matches, err := filepath.Glob(tpath)
	if err != nil {
		panic(err)
	}

	var tasks []Task
	for _, match := range matches {
		task := NewTask(match)

		switch task.Config.Setting.Mode {
		case ModeCookie:
			tasks = append(tasks, *task)
		case ModeProxy:
			for _, t := range task.SplitFromProxies() {
				tasks = append(tasks, t)
			}
		}
	}
	p.Tasks = tasks

	return &p
}

func GoResponseToPy(gresp *requests.Response) *py.PyObject {
	obj := py.PyDict_New()
	py.PyDict_SetItem(obj, py.PyString_FromString("Status"), py.PyInt_FromLong(gresp.GResponse.StatusCode))
	py.PyDict_SetItem(obj, py.PyString_FromString("Content"), py.PyString_FromString(gresp.Content()))
	return obj
}

// Execute 执行
func (person *Person) Execute() {
	//TODO: Python的脚本函数, 与动态更新 返回的数据在python的脚本处理
	for _, task := range person.Tasks {
		for _, PResult := range task.ExecuteOnPlan() {
			if task.Config.Setting.ResultProcessing != "" {
				GoResponseToPy(PResult.Resp)
			}
		}

	}
}
