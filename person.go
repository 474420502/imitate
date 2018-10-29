package main

// Person 拥有两种类型并发任务状态, 一种是以cookie为个例控制, 另一种是以 代理为个例子.
type Person struct {
}

// NewPerson 创建一个新Person
func NewPerson() *Person {
	p := Person{}
	return &p
}
