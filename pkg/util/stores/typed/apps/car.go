package apps

import (
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
	List() (appleList *appsv1.CarList, err error)
	Create(apple *appsv1.Car) (ver *appsv1.Car, err error)
	Delete(name string) (err error)
	Update(apple *appsv1.Car) (ver *appsv1.Car, err error)

}

type car struct {
	client rest.Interface
}

// Get 获取apple资源
func (v *car) Get(name string) (ver *appsv1.Car, err error) {
	ver = &appsv1.Car{}
	err = v.client.
		Get().
		Path("/car").GetCarName(name).
		Do().
		Into(ver)

	return
}

// Post 创建apple资源
func (v *car) Create(car *appsv1.Car) (ver *appsv1.Car, err error) {
	ver = &appsv1.Car{}
	err = v.client.
		Post().Path("/car").CreateCar(car).
		Do().Into(ver)
	return
}

func (v *car) List() (carList *appsv1.CarList, err error) {
	carList = &appsv1.CarList{}
	err = v.client.Get().Path("/carlist").ListCar().Do().Into(carList)

	return
}

func (v *car) Delete(name string) (err error) {
	ver := &appsv1.Car{}
	err = v.client.Delete().Path("/car").DeleteCar(name).Do().Into(ver)

	return
}

func (v *car) Update(car *appsv1.Car) (ver *appsv1.Car, err error) {
	ver = &appsv1.Car{}
	err = v.client.
		Put().Path("/car").UpdateCar(car).
		Do().Into(ver)
	return
}





var _ CarInterface = &car{}

