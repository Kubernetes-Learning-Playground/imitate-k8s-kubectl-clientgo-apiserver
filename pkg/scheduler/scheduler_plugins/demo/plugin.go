package demo

import (
	"math/rand"
	"practice_ctl/pkg/scheduler/nodes"
	"practice_ctl/pkg/scheduler/scheduler_plugins"
)

type MockPlugin struct{}

func (m MockPlugin) Score(task scheduler_plugins.SchedulerObject, infos *nodes.NodeInfo) float64 {
	return rand.Float64() * 100
}

func (m MockPlugin) Filter(task scheduler_plugins.SchedulerObject) bool {
	return rand.Float64() != 0
}

var _ scheduler_plugins.Plugin = &MockPlugin{}
