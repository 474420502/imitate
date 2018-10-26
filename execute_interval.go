package main

import (
	"time"

	py "github.com/sbinet/go-python"
)

// ExecuteInterval 时间间隔的类型
type ExecuteInterval struct {
	TimeInterval int64 // 时间间隔

	TriggerTime int64 // 执行时间间隔触发时间
	StartStatus bool  // 判断是否按照时间间隔执行
}

// SetStartStatus 设置执行计划是否生效
func (ei *ExecuteInterval) SetStartStatus(status bool) {
	ei.StartStatus = status
}

// GetTriggerTime 获取计划的触发时间
func (ei *ExecuteInterval) GetTriggerTime() int64 {
	if ei.StartStatus {
		return ei.TriggerTime
	}
	return -1
}

// TimeTo 是否到了该触发的时间
func (ei *ExecuteInterval) TimeTo() int64 {
	return time.Now().Unix() - ei.TriggerTime
}

// GetStartStatus 获取计划的触发时间是否在生效
func (ei *ExecuteInterval) GetStartStatus() bool {
	return ei.StartStatus
}

// CalculateTrigger 计算触发特定时间任务的时间点
func (ei *ExecuteInterval) CalculateTrigger() int64 {
	now := time.Now()
	ei.TriggerTime = now.Unix() + ei.TimeInterval
	return ei.TriggerTime
}

// FromPyObject 从python object里获取执行时间间隔
func (ei *ExecuteInterval) FromPyObject(obj *py.PyObject) {
	ei.FromValue(py.PyLong_AsLong(obj))
}

// FromValue 生成计划表
func (ei *ExecuteInterval) FromValue(vsleep int64) {
	ei.TimeInterval = vsleep
}
