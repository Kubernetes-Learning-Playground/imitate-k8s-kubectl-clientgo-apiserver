package stores

import (
	"practice_ctl/pkg/util/stores/rest"
	"practice_ctl/pkg/util/stores/typed/core"
)

// 模仿k8s 的clientset
// 调用方式：clientSet.Core().Version().Get()

// ClientSet 客户端
type ClientSet struct {
	*rest.RESTClient
}

func (cs *ClientSet) Core() core.CoreInterface {
	return core.New(cs.RESTClient)
}

// NewForConfig 读配置文件
func NewForConfig(c *rest.Config) *ClientSet {
	rc := rest.NewRESTClient(c)
	return &ClientSet{RESTClient: rc}
}
