package scheduler_plugins

import (
	"practice_ctl/pkg/apimachinery/runtime"
)

type schedulerObject struct {
	Object           runtime.Object
	Name             string
	SelectNode       string
	ObjectRecordNode []ObjectRecordNode
}

// ObjectRecordNode 为了记录特定object在所有node中的打分
type ObjectRecordNode struct {
	NodeName string
	Score    float64
}

func (m *schedulerObject) SetNode(s string) {
	// FIXME 这里需要改，需要在runtime.Object中修改SelectNode 不是在外层
	m.SelectNode = s
}

func (m *schedulerObject) SetObjectRecordNode(s string, f float64) {
	r := ObjectRecordNode{
		Score:    f,
		NodeName: s,
	}
	m.ObjectRecordNode = append(m.ObjectRecordNode, r)
}

func (m *schedulerObject) GetObjectRecordNodeList() []ObjectRecordNode {
	return m.ObjectRecordNode
}

var _ SchedulerObject = &schedulerObject{}
