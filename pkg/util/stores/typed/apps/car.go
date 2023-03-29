package apps

import (
	"practice_ctl/pkg/apimachinery/runtime"
	appsv1 "practice_ctl/pkg/apis/apps/v1"
	"practice_ctl/pkg/util/stores/rest"
)

type CarRequest struct {
	*rest.Request
	car *appsv1.Car
}

type CarGetter interface {
	Car() CarInterface
}

func newCar(c rest.Interface) CarInterface {
	return &car{client: c}
}


type CarInterface interface {
	Get(name string) (ver *appsv1.Car, err error)
	List() (carList *appsv1.CarList, err error)
	Create(car runtime.Object) (ver *appsv1.Car, err error)
	Delete(name string) (err error)
	Update(car runtime.Object) (ver *appsv1.Car, err error)
	Watch() *rest.Request
}

type car struct {
	client rest.Interface
}

// Get 获取apple资源
func (v *car) Get(name string) (ver *appsv1.Car, err error) {
	ver = &appsv1.Car{}
	err = v.client.
		Get().
		Path("/v1/car").GetCarName(name).
		Do().
		Into(ver)

	return
}

// Post 创建apple资源
func (v *car) Create(car runtime.Object) (ver *appsv1.Car, err error) {
	ver = &appsv1.Car{}
	err = v.client.
		Post().Path("/v1/car").CreateCar(car.(*appsv1.Car)).
		Do().Into(ver)
	return
}

func (v *car) List() (carList *appsv1.CarList, err error) {
	carList = &appsv1.CarList{}
	err = v.client.Get().Path("/v1/carlist").ListCar().Do().Into(carList)

	return
}

func (v *car) Delete(name string) (err error) {
	ver := &appsv1.Car{}
	err = v.client.Delete().Path("/v1/car").DeleteCar(name).Do().Into(ver)

	return
}

func (v *car) Update(car runtime.Object) (ver *appsv1.Car, err error) {
	ver = &appsv1.Car{}
	err = v.client.
		Put().Path("/v1/car").UpdateCar(car.(*appsv1.Car)).
		Do().Into(ver)
	return
}


func (v *car) Watch() *rest.Request {

	res := v.client.
		Watch().WsPath("/v1/car/watch").
		WatchCar()

	return res
}


var _ CarInterface = &car{}

