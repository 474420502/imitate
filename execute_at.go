package main

import (
	"time"

	py "github.com/sbinet/go-python"
)

// ExecuteAt 特定的时间任务 接口源自于 IExecute
type ExecuteAt struct {
	Year  int
	Month int
	Day   int
	Hour  int
	Min   int
	Sec   int

	TriggerTime int64 // 下次的触发时间点
	StartStatus bool  // 一个值判断这个时间表是否有效
}

// SetStartStatus 设置执行计划是否生效
func (ea *ExecuteAt) SetStartStatus(status bool) {
	ea.StartStatus = status
}

// GetTriggerTime 获取计划的触发时间
func (ea *ExecuteAt) GetTriggerTime() int64 {
	return ea.TriggerTime
}

// GetStartStatus 获取计划的触发时间是否在生效
func (ea *ExecuteAt) GetStartStatus() bool {
	return ea.StartStatus
}

// CalculateTrigger 计算触发特定时间任务的时间点 执行后 可以通过GetTriggerTime确认触发时间
func (ea *ExecuteAt) CalculateTrigger() int64 {
	now := time.Now()

	year := ea.Year
	if ea.Year <= 0 {
		year = now.Year()
	}

	month := time.Month(ea.Month)
	if ea.Month <= 0 {
		month = now.Month()
	}

	day := ea.Day
	if ea.Day <= 0 {
		day = now.Day()
	}

	hour := ea.Hour
	if ea.Hour < 0 {
		hour = now.Hour()
	}

	min := ea.Min
	if ea.Min < 0 {
		min = now.Minute()
	}

	sec := ea.Sec
	if ea.Sec < 0 {
		sec = now.Second()
	}

	ea.TriggerTime = time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local).Unix()
	return ea.TriggerTime
}

// FromPyObject 从python object里获取执行时间表的结构 object 为size:6 的tuple
func (ea *ExecuteAt) FromPyObject(obj *py.PyObject) {

	ea.FromValues(
		py.PyInt_AsLong(py.PyTuple_GetItem(obj, 0)),
		py.PyInt_AsLong(py.PyTuple_GetItem(obj, 1)),
		py.PyInt_AsLong(py.PyTuple_GetItem(obj, 2)),
		py.PyInt_AsLong(py.PyTuple_GetItem(obj, 3)),
		py.PyInt_AsLong(py.PyTuple_GetItem(obj, 4)),
		py.PyInt_AsLong(py.PyTuple_GetItem(obj, 5)),
	)
}

// FromValues 从数值 里获取执行时间表的结构
func (ea *ExecuteAt) FromValues(year int, month int, day int, hour int, min int, sec int) {
	ea.Year = year
	ea.Month = month
	ea.Day = day
	ea.Hour = hour
	ea.Min = min
	ea.Sec = sec
}
