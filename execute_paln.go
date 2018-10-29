package main

import (
	"github.com/474420502/requests"
)

// type ExecuteRecord struct {
// 	IsSuccess bool
// 	FailCount int
// 	Message   []string
// }

// IExecute 计划执行的时间接口
type IExecute interface {
	SetStartStatus(status bool) // SetStart 设置执行计划是否生效
	GetStartStatus() bool       // IsStart 获取计划的触发时间是否在生效
	GetTriggerTime() int64      // GetTriggerTime 获取计划的触发时间

	TimeTo() int64 // TimeTo 是否到了该触发的时间

	// SetSuccessStatus(status bool) // SetSuccessStatus 设置成功的状态 将记录于历史
	// History() []ExecuteRecord     // History 记录一些历史, 可能会持久到数据库. 暂时不要
	CalculateTrigger() int64 // CalculateTrigger 计算触发特定时间任务的时间点
}

// ExecutePlan 执行时间的计划表
type ExecutePlan struct {
	ExecuteQueue []IExecute
}

// PlanResult 执行计划表的结果
type PlanResult struct {
	Exec IExecute
	Resp *requests.Response
}

// AppendIExecute 添加执行计划任务
func (ep *ExecutePlan) AppendIExecute(e IExecute) {
	ep.ExecuteQueue = append(ep.ExecuteQueue, e)
}

// ClearIExecute 清除执行计划任务
func (ep *ExecutePlan) ClearIExecute() {
	ep.ExecuteQueue = []IExecute{}
}

// CountIExecute 清除执行计划任务
func (ep *ExecutePlan) CountIExecute() int {
	return len(ep.ExecuteQueue)
}
