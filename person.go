package main

import "github.com/474420502/grequests"

// Person 爬虫个例 以Cookie为个例 或者 以代理IP为个例
type Person struct {
	Config  *TaskConfig
	Session *grequests.Session
}

// NewPerson new 一个person 对象
func NewPerson(taskFileName string) *Person {

	p := &Person{}
	p.Config = &TaskConfig{}
	p.Config.Load(taskFileName)

	p.Session = grequests.NewSession(nil)

	return p
}
