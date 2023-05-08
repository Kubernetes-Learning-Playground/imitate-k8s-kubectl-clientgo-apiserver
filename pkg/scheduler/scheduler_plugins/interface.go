package scheduler_plugins

import (
	"practice_ctl/pkg/scheduler/nodes"
)

// Pod 接口
type Pod interface {
	Exec()
	SetNode(string)
	SetPodRecordNode(string, float64)
	GetPodRecordNodeList() []PodRecordNode
}

type SchedulerObject interface {
	SetNode(string)
	SetObjectRecordNode(string, float64)
	GetObjectRecordNodeList() []ObjectRecordNode
}

// PodRecordNode 为了记录特定pod在所有node中的打分
type PodRecordNode struct {
	NodeName string
	Score    float64
}

// Plugin 插件接口
type Plugin interface {
	Score(pod SchedulerObject, infos *nodes.NodeInfo) float64
	Filter(pod SchedulerObject) bool
}
