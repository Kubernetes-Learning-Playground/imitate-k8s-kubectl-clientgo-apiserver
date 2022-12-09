package core

import (
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

// 本课程来自 程序员在囧途(www.jtthink.com) 咨询群：98514334
type AppleInterface interface {
	Get(name string) (ver *v1.Apple, err error)
	List() (appleList *v1.AppleList, err error)
	Create(apple *v1.Apple) (ver *v1.Apple, err error)
	Delete(name string) (err error)
	Update(apple *v1.Apple) (ver *v1.Apple, err error)

}

type apple struct {
	client rest.Interface
}

// Get 获取apple资源
func (v *apple) Get(name string) (ver *v1.Apple, err error) {
	ver = &v1.Apple{}
	err = v.client.
		Get().
		Path("/apple").GetAppleName(name).
		Do().
		Into(ver)

	return
}

// Post 创建apple资源
func (v *apple) Create(apple *v1.Apple) (ver *v1.Apple, err error) {
	ver = &v1.Apple{}
	err = v.client.
		Post().Path("/apple").CreateApple(apple).
		Do().Into(ver)
	return
}

func (v *apple) List() (appleList *v1.AppleList, err error) {
	appleList = &v1.AppleList{}
	err = v.client.Get().Path("/applelist").ListApple().Do().Into(appleList)

	return
}

func (v *apple) Delete(name string) (err error) {
	ver := &v1.Apple{}
	err = v.client.Delete().Path("/apple").DeleteApple(name).Do().Into(ver)

	return
}

func (v *apple) Update(apple *v1.Apple) (ver *v1.Apple, err error) {
	ver = &v1.Apple{}
	err = v.client.
		Put().Path("/apple").UpdateApple(apple).
		Do().Into(ver)
	return
}





var _ AppleInterface = &apple{}
