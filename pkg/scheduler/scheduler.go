package scheduler

import (
	"context"
	"sync"
)

type Scheduler struct {

	queue           *taskQueue                       // 任务队列
	workerChan      chan struct{}                    // 控制并发数用的channel
	//responseWriter  func(result *model.Result) error // 向Server响应结果用
	workerWaitGroup *sync.WaitGroup
	closed          bool

}

func New() (*Scheduler, error) {
	return &Scheduler{}, nil
}


func (s *Scheduler) Run(ctx context.Context) {

}