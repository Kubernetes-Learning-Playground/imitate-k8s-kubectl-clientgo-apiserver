package core

import (
	"practice_ctl/pkg/apimachinery/runtime"
	v1 "practice_ctl/pkg/apis/core/v1"
	"practice_ctl/pkg/util/stores/rest"
)

type AppleRequest struct {
	*rest.Request
	apple *v1.Apple
}

type AppleGetter interface {
	Apple() AppleInterface
}

func newApple(c rest.Interface) AppleInterface {
	return &apple{client: c}
}

type AppleInterface interface {
	Get(name string) (ver *v1.Apple, err error)
	List() (appleList *v1.AppleList, err error)
	Create(apple runtime.Object) (ver *v1.Apple, err error)
	Delete(name string) (err error)
	Update(apple runtime.Object) (ver *v1.Apple, err error)
	Watch() *rest.Request
}

type apple struct {
	client rest.Interface
}

// Get 获取apple资源
func (v *apple) Get(name string) (ver *v1.Apple, err error) {
	ver = &v1.Apple{}
	err = v.client.
		Get().
		Path("/v1/apple").GetAppleName(name).
		Do().
		Into(ver)

	return
}

// Post 创建apple资源
func (v *apple) Create(apple runtime.Object) (ver *v1.Apple, err error) {
	ver = &v1.Apple{}
	err = v.client.
		Post().Path("/v1/apple").CreateApple(apple.(*v1.Apple)).
		Do().Into(ver)
	return
}

func (v *apple) List() (appleList *v1.AppleList, err error) {
	appleList = &v1.AppleList{}
	err = v.client.Get().Path("/v1/applelist").ListApple().Do().Into(appleList)

	return
}

func (v *apple) Delete(name string) (err error) {
	ver := &v1.Apple{}
	err = v.client.Delete().Path("/v1/apple").DeleteApple(name).Do().Into(ver)

	return
}

func (v *apple) Update(apple runtime.Object) (ver *v1.Apple, err error) {
	ver = &v1.Apple{}

	err = v.client.
		Put().Path("/v1/apple").UpdateApple(apple.(*v1.Apple)).
		Do().Into(ver)
	return
}

func (v *apple) Watch() *rest.Request {

	res := v.client.
		Watch().WsPath("/v1/apple/watch").
		WatchApple()

	return res
}




var _ AppleInterface = &apple{}
