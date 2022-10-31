package blogs

import (
	"practice_ctl/pkg/util/blogs/rest"
	"practice_ctl/pkg/util/blogs/typed/core"
)

// 本课程来自 程序员在囧途(www.jtthink.com) 咨询群：98514334
// 模仿k8s 的clientset
type Clientset struct {
	*rest.RESTClient
}

func (cs *Clientset) Core() core.CoreInterface {
	return core.New(cs.RESTClient)
}
func NewForConfig(c *rest.Config) *Clientset {
	rc := rest.NewRESTClient(c)
	return &Clientset{RESTClient: rc}
}
