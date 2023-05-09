package scheduler

import (
	"context"
	"fmt"
	"k8s.io/klog/v2"
	"practice_ctl/pkg/scheduler/nodes"
	"practice_ctl/pkg/scheduler/queue"
	"practice_ctl/pkg/scheduler/scheduler_plugins"
	"practice_ctl/pkg/scheduler/scheduler_plugins/demo"
	"sync"
)

type Scheduler struct {
	options *schedulerOptions // 调度器配置
	name    string

	queue     *queue.Queue
	pods      chan scheduler_plugins.SchedulerObject // pod队列
	nodeInfos *nodes.NodeInfos                       // 存储所有node的信息
	workers   int                                    // 控制并发数
	plugins   []scheduler_plugins.Plugin             // 插件

	wg     sync.WaitGroup
	stopC  chan struct{} // 通知
	logger klog.Logger
}

type schedulerOptions struct {
	healthPort    int
	numWorker     int
	queueCapacity int
}

// defaultOptions 默认配置
var defaultOptions = schedulerOptions{
	numWorker:     3,
	queueCapacity: 10,
	healthPort:    9999,
}

func New() (*Scheduler, error) {
	schedulerName := fmt.Sprintf("%s-scheduler", "store")

	// 调度器
	scheduler := Scheduler{
		name:      schedulerName,
		options:   &defaultOptions,
		queue:     queue.NewScheduleQueue(defaultOptions.queueCapacity),
		nodeInfos: nodes.NewNodeInfos(),
		workers:   defaultOptions.numWorker,
		plugins:   make([]scheduler_plugins.Plugin, 0),
		wg:        sync.WaitGroup{},
		stopC:     make(chan struct{}),
		logger:    klog.LoggerWithName(klog.Background(), schedulerName),
	}

	// 加入模拟node
	scheduler.nodeInfos.AddNode(nodes.NewNodeInfo("node1"))
	scheduler.nodeInfos.AddNode(nodes.NewNodeInfo("node2"))
	scheduler.nodeInfos.AddNode(nodes.NewNodeInfo("node3"))

	// 加入plugin插件
	// FIXME:
	scheduler.AddPlugin(&demo.MockPlugin{})

	return &scheduler, nil
}

func (s *Scheduler) Run(ctx context.Context) {
	// 启动队列
	go s.queue.Run(s.stopC)

	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		// 并发
		go func() {
			defer s.wg.Done()

			for {
				select {
				case <-s.stopC: // 退出通知
					s.logger.Info("return scheduler...")
					return
				case t := <-s.queue.Get(): // 取出队列中的pod

					t = s.runFilter(t)
					// 打分node
					t = s.runScorer(t)
					// 选出最好的node
					nodeInfo := s.selectHost(t)
					// 异步绑定
					go bind(t, nodeInfo)
				}
			}
		}()
	}

	select {
	case <-ctx.Done():
		klog.Info("scheduler stop...")
		return
	}
}

// Stop 停止
func (s *Scheduler) Stop() {
	if s.queue.Len() > 0 {
		s.logger.Info("scheduler queue still have element...")
	}
	// 通知退出
	close(s.stopC) // 通知
}

// AddPlugin 加入插件
func (s *Scheduler) AddPlugin(plugin scheduler_plugins.Plugin) {
	s.plugins = append(s.plugins, plugin)
}

// AddPod 放入pod
func (s *Scheduler) AddPod(pod scheduler_plugins.SchedulerObject) {
	s.queue.ActiveQ <- pod
}

// runFilter 执行过滤插件
func (s *Scheduler) runFilter(pod scheduler_plugins.SchedulerObject) scheduler_plugins.SchedulerObject {
	s.logger.Info("runFilter...")
	if s.podFiltered(pod) {
		return pod
	}
	s.logger.Info("have no pod to run...")
	// 把无法调度的放入backoffQ中
	s.queue.Backoff(pod)
	return nil
}

// runScorer 执行打分插件
func (s *Scheduler) runScorer(pod scheduler_plugins.SchedulerObject) scheduler_plugins.SchedulerObject {
	s.logger.Info("runScorer...")
	var totalScore float64
	if pod == nil {
		s.logger.Info("have no pod to score...")
		return nil
	}
	for _, nodeInfo := range s.nodeInfos.NodeInfos {
		for _, plugin := range s.plugins {
			totalScore += plugin.Score(pod, nodeInfo)
		}
		if totalScore != 0 {
			pod.SetObjectRecordNode(nodeInfo.NodeName, totalScore)
			totalScore = 0
		}
	}

	return pod
}

// selectHost 选出最好的node
func (s *Scheduler) selectHost(pod scheduler_plugins.SchedulerObject) *nodes.NodeInfo {
	var resNodeInfo *nodes.NodeInfo
	if pod == nil {
		s.logger.Info("have no pod to select host...")
		return resNodeInfo
	}
	nodeList := pod.GetObjectRecordNodeList()
	var maxScore float64
	var maxScoreNodeName string
	for _, node := range nodeList {
		if node.Score >= maxScore {
			maxScore = node.Score
			maxScoreNodeName = node.NodeName
		}
	}

	for _, schedulerNode := range s.nodeInfos.NodeInfos {
		if schedulerNode.NodeName == maxScoreNodeName {
			resNodeInfo = schedulerNode
		}
	}
	return resNodeInfo
}

// bind 绑定
func bind(pod scheduler_plugins.SchedulerObject, nodeInfo *nodes.NodeInfo) {
	if pod == nil || nodeInfo == nil {
		klog.Info("have no pod or node to bind...")
		return
	}
	pod.SetNode(nodeInfo.NodeName)
}

// podFiltered 过滤插件
func (s *Scheduler) podFiltered(pod scheduler_plugins.SchedulerObject) bool {
	for _, plugin := range s.plugins {
		if plugin.Filter(pod) {
			return true
		}
	}
	return false
}
