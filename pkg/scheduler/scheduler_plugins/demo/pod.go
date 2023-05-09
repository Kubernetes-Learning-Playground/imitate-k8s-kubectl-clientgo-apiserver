package demo

import (
	"k8s.io/klog/v2"
	"practice_ctl/pkg/scheduler/scheduler_plugins"
)

type MockPod struct {
	Name          string
	SelectNode    string
	PodRecordNode []scheduler_plugins.PodRecordNode
}

func (m *MockPod) Exec() {
	klog.Info("pod ", m.Name, " in ", m.SelectNode, " node")
}

func (m *MockPod) SetNode(s string) {
	m.SelectNode = s
}

func (m *MockPod) SetPodRecordNode(s string, f float64) {
	r := scheduler_plugins.PodRecordNode{
		Score:    f,
		NodeName: s,
	}
	m.PodRecordNode = append(m.PodRecordNode, r)
}

func (m *MockPod) GetPodRecordNodeList() []scheduler_plugins.PodRecordNode {
	return m.PodRecordNode
}

var _ scheduler_plugins.Pod = &MockPod{}
