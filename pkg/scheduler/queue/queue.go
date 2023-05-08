package queue

import "practice_ctl/pkg/scheduler/scheduler_plugins"

// Queue 调度队列
type Queue struct {
	// 所有加入队列的对象都放入此chan
	ActiveQ chan scheduler_plugins.SchedulerObject
	// 当有调度错误时，放入backoffQ
	backoffQ chan scheduler_plugins.SchedulerObject

	// 调度时，从out chan中取对象
	out chan scheduler_plugins.SchedulerObject
}

func NewScheduleQueue(capacity int) *Queue {
	return &Queue{
		ActiveQ:  make(chan scheduler_plugins.SchedulerObject, capacity),
		backoffQ: make(chan scheduler_plugins.SchedulerObject, capacity),
		out:      make(chan scheduler_plugins.SchedulerObject, capacity*2),
	}
}

// Run 放入out chan的逻辑
// 一般情况：如果没有调度失败，就会把scheduleQ中的对象放入out中，
// 如果调度失败，就有概率从backoffQ中获取对象放入out中
func (q *Queue) Run(done <-chan struct{}) {
	for {
		select {
		case <-done:
			return
		case t := <-q.backoffQ:
			q.out <- t
		default:
			select {
			case <-done:
				return
			case t := <-q.backoffQ:
				q.out <- t
			case t := <-q.ActiveQ:
				q.out <- t
			}
		}
	}
}

// Get first try backoffQ, and then scheduleQ. It blocks if both backoffQ and scheduleQ are empty
func (q *Queue) Get() <-chan scheduler_plugins.SchedulerObject {
	return q.out
}

// Put 入队
func (q *Queue) Put(t scheduler_plugins.SchedulerObject) {
	q.ActiveQ <- t
}

// Backoff 调度错误时，入队
func (q *Queue) Backoff(t scheduler_plugins.SchedulerObject) {
	q.backoffQ <- t
}

func (q *Queue) Len() int {
	return len(q.ActiveQ) + len(q.backoffQ)
}