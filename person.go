package main

import "path/filepath"

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
	mode  TypeMode
	Tasks []Task
}

// NewPerson 创建一个新Person
func NewPerson() *Person {
	p := Person{}

	matches, err := filepath.Glob("task/*_config.py")
	if err != nil {
		panic(err)
	}

	var tasks []*Task
	for _, match := range matches {
		task := NewTask(match)
		tasks = append(tasks, task)
	}

	switch m {
	case ModeCookie:

	case ModeProxy:

	}

	return &p
}

// SetMode 设置任务的模式
// func (p *Person) SetMode(m TypeMode) {
// 	p.mode = m
// }

//GetMode 获取任务的模式
func (p *Person) GetMode(m TypeMode) TypeMode {
	return p.mode
}
