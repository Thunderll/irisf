package pub_sub

import (
	"iris_project_foundation/config"
	"time"
)

const (
	Free      = iota // 空闲
	Executing        // 执行中
	Finished         // 结束
)

// ScheduleTask 调度任务
type ScheduleTask struct {
	OrderSN string `json:"orderSN"` // 订单号
	LockKey string `json:"lockKey"` // 数据锁
	Status  int    `json:"status"`  // 任务状态, 0-未执行, 1-执行中, 2-完成
}

func BuildScheduleTask(orderSN string) *ScheduleTask {
	return &ScheduleTask{
		OrderSN: orderSN,
		LockKey: config.GConfig.App.OrderLockPrefix + orderSN,
		Status:  Free,
	}
}

// ScheduleResult 任务执行结果
type ScheduleResult struct {
	Task    *ScheduleTask
	Success bool
	Err     error
	EndTime time.Time
}

func BuildScheduleResult(task *ScheduleTask, success bool, err error) *ScheduleResult {
	return &ScheduleResult{
		Task:    task,
		Success: success,
		Err:     err,
		EndTime: time.Now(),
	}
}
