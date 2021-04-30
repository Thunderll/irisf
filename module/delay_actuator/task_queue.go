package pub_sub

import (
	"context"
	"encoding/json"
	"iris_project_foundation/module/redis_manager"

	"github.com/go-redis/redis/v8"
)

type ITaskQueue interface {
	Add(string) error
	Remove(*ScheduleResult) error
	Queue() chan *ScheduleTask
}

// TaskQueue 任务队列, 借助redis作为本地任务的备份, 防止宕机丢失
type TaskQueue struct {
	queue chan *ScheduleTask
}

// BuildTaskQueue 初始化任务队列
// 初始化时从redis中同步到本地
func BuildTaskQueue() *TaskQueue {
	var (
		err   error
		res   *redis.StringStringMapCmd
		queue chan *ScheduleTask
		task  *ScheduleTask
	)
	queue = make(chan *ScheduleTask, 1000)
	res = redis_manager.RDB.HGetAll(context.TODO(), "taskQueue")
	for _, v := range res.Val() {
		if err = json.Unmarshal([]byte(v), task); err != nil {
			// TODO log记录一下从redis取出失败的任务
			continue
		}
		queue <- task
	}
	return &TaskQueue{queue: queue}
}

func (q *TaskQueue) Add(orderSN string) (err error) {
	var (
		task     *ScheduleTask
		taskJson []byte
	)

	task = BuildScheduleTask(orderSN)
	if taskJson, err = json.Marshal(task); err != nil {
		return err
	}
	// 将任务备份到redis中
	redis_manager.RDB.HSetNX(context.TODO(), "taskQueue", orderSN, taskJson)
	// 加入本地任务队列
	q.queue <- task

	return
}

func (q *TaskQueue) Remove(result *ScheduleResult) (err error) {
	// 将任务从redis中删除
	redis_manager.RDB.HDel(context.TODO(), "taskQueue", result.Task.OrderSN)

	return
}

func (q *TaskQueue) Queue() chan *ScheduleTask {
	return q.queue
}
