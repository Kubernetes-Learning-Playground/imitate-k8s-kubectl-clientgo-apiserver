package core

import (
	"practice_ctl/pkg/apis/core/unverstioned"
	"practice_ctl/pkg/util/blogs/rest"
)

type VersionGetter interface {
	Version() VersionInterface
}

func newVersion(c rest.Interface) VersionInterface {
	return &version{client: c}
}

// 本课程来自 程序员在囧途(www.jtthink.com) 咨询群：98514334
type VersionInterface interface {
	Get() (ver *unverstioned.Version, err error)
}

type version struct {
	client rest.Interface
}

func (v *version) Get() (ver *unverstioned.Version, err error) {
	//TODO implement me
	ver = &unverstioned.Version{}
	err = v.client.
		Get().
		Path("/version").
		Do().
		Into(ver)
	return
}

var _ VersionInterface = &version{}
