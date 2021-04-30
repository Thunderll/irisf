package pub_sub

import (
	"context"
	rdb "iris_project_foundation/module/redis_manager"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var GScheduler *Scheduler

type Scheduler struct {
	EventChan  <-chan *redis.Message // 订单过期事件队列
	TaskQueue  ITaskQueue            // 任务调度队列
	ResultChan chan *ScheduleResult  // 任务执行结果队列
}

func InitScheduler() (err error) {
	GScheduler = &Scheduler{
		EventChan:  rdb.RDB.Subscribe(context.TODO(), "__keyevent@0__:expired").Channel(),
		TaskQueue:  BuildTaskQueue(),
		ResultChan: make(chan *ScheduleResult, 1000),
	}

	go GScheduler.ScheduleLoop()
	return nil
}

// ScheduleLoop 调度循环
func (s *Scheduler) ScheduleLoop() {
	var (
		event  *redis.Message
		result *ScheduleResult
		task   *ScheduleTask
	)

	for {
		select {
		case event = <-s.EventChan: // 有订单到期事件
			s.HandleEvent(event)
		case result = <-s.ResultChan: // 任务执行完成
			s.HandleResult(result)
		case task = <-s.TaskQueue.Queue(): // 执行调度
			s.TrySchedule(task)
		}
	}
}

// PushTask 将订单加入延迟执行队列中
func PushTask(orderSN string, expiration int64) {
	log.Printf("[INFO] 导入订单(%s)", orderSN)
	// 将订单好放入redis中并设置过期时间
	rdb.RDB.SetNX(context.TODO(), orderSN, 1, time.Duration(expiration)*time.Second)
}

// HandleEvent 提取订单到期事件, 放入任务队列
func (s *Scheduler) HandleEvent(event *redis.Message) {
	if strings.HasPrefix(event.Payload, "orderTask") {
		_ = s.TaskQueue.Add(strings.TrimLeft(event.Payload, "orderTask:"))
	}
}

// TrySchedule 尝试调度任务
func (s *Scheduler) TrySchedule(task *ScheduleTask) {
	var (
		ok bool
	)
	// 尝试上锁, 上锁失败暂时跳过
	if ok = rdb.LockKey(task.LockKey); !ok || task.Status != Free {
		return
	}
	s.ExecuteTask(task)
}

// ExecuteTask 执行任务
func (s *Scheduler) ExecuteTask(task *ScheduleTask) {
	task.Status = Executing
	// 模拟处理任务
	go func() {
		log.Printf("[DEBUG] 执行任务(%s)...", task.OrderSN)
		time.Sleep(2 * time.Second)

		PushResult(BuildScheduleResult(task, true, nil))
		log.Printf("[DEBUG] 任务(%s)执行完毕.", task.OrderSN)
	}()
}

// PushResult 将任务结果放入结果队列中
func PushResult(result *ScheduleResult) {
	result.Task.Status = Finished
	GScheduler.ResultChan <- result
}

// HandleResult 处理任务结果
func (s *Scheduler) HandleResult(result *ScheduleResult) {
	if !result.Success {
		log.Printf("[ERR] 任务(%s)执行失败: %v", result.Task.OrderSN, result.Err)
	} else {
		log.Printf("[INFO] 任务(%s)执行成功.", result.Task.OrderSN)
	}

	// 处理结果
	log.Println("[DEBUG] 处理执行结果.")
	time.Sleep(1 * time.Second)

	// 从任务队列中删除, 实际删除redis中的备份
	_ = s.TaskQueue.Remove(result)
	rdb.Unlock(result.Task.LockKey)
}
