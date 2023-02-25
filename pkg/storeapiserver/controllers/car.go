package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s.io/klog/v2"
	appsv1 "practice_ctl/pkg/apis/apps/v1"
	"practice_ctl/pkg/etcd"
)

var CarMap = map[string]*appsv1.Car{}

func init() {
	// 初始化一个对象
	init := &appsv1.Car{
		ApiVersion: "apps/v1",
		Kind: 		"Car",
		Metadata: appsv1.Metadata{
			Name: "initCar",
		},
		Spec: appsv1.CarSpec{
			Brand: "initBrand",
			Price: "initPrice",
			Color: "initColor",
		},
		Status: appsv1.CarStatus{
			Status: "created",
		},
	}
	CarMap[init.Name] = init

}


func getCar(name string) (*appsv1.Car, error) {
	if car, ok := CarMap[name]; ok {
		return car, nil
	}
	return nil, nil
}

func deleteCar(name string) error {
	if car, ok := CarMap[name]; ok {
		strKey, _ := parseEtcdDataCar(car)
		klog.Info("delete key: ", strKey)
		delete(CarMap, car.Name)
		err := etcd.Delete(strKey)
		if err != nil {
			klog.Errorf("delete key error: ", strKey, err)
			return err
		}
		return nil
	}
	return nil

}

func listCar() (*appsv1.CarList, error) {
	carList := &appsv1.CarList{}
	for _, v := range CarMap {
		carList.Item = append(carList.Item, v)
	}

	return carList, nil
}

func parseEtcdDataCar(car *appsv1.Car) (string, string) {
	strKey := "/" + car.Kind + "/" + car.Name
	strValue := car.Name

	return strKey, strValue
}

func createOrUpdateCar(car *appsv1.Car) (*appsv1.Car, error) {
	if old, ok := CarMap[car.Name]; ok {
		klog.Infof("find the apple %v, and update it!", old.Name)
		old.Name = car.Name
		old.Spec.Price = car.Spec.Price
		old.Spec.Brand = car.Spec.Brand
		old.Spec.Color = car.Spec.Color
		old.Status.Status = "updated"

		strKey, strValue := parseEtcdDataCar(car)
		klog.Info("update key: ", strKey)
		err := etcd.Put(strKey, strValue)
		if err != nil {
			klog.Errorf("update key error: ", strKey, err)
			return old, err
		}
		return old, nil
	}
	klog.Infof("not find this car, and create it!")

	new := &appsv1.Car{
		ApiVersion: "apps/v1",
		Kind: "Car",
		Metadata: appsv1.Metadata{
			Name: car.Name,
		},
		Spec: appsv1.CarSpec{
			Brand: car.Spec.Brand,
			Price: car.Spec.Price,
			Color: car.Spec.Color,
		},
		Status: appsv1.CarStatus{
			Status: "created",
		},
	}


	// 存入map
	CarMap[car.Name] = new

	return new, nil
}


func createCar(car *appsv1.Car) (*appsv1.Car, error) {
	// 如果查到就抛错
	if _, ok := CarMap[car.Name]; ok {
		return nil, errors.New("this car is created ")
	}
	new := &appsv1.Car{
		ApiVersion: "apps/v1",
		Kind: "Car",
		Metadata: appsv1.Metadata{
			Name: car.Name,
		},
		Spec: appsv1.CarSpec{
			Brand: car.Spec.Brand,
			Price: car.Spec.Price,
			Color: car.Spec.Color,
		},
		Status: appsv1.CarStatus{
			Status: "created",
		},
	}


	// 存入map
	CarMap[car.Name] = new

	return new, nil

}

func updateCar(car *appsv1.Car) (*appsv1.Car, error) {
	// 重新赋值
	if old, ok := CarMap[car.Name]; ok {
		old.Name = car.Name
		old.Spec.Price = car.Spec.Price
		old.Spec.Brand = car.Spec.Brand
		old.Spec.Color = car.Spec.Color
		return old, nil
	}


	// 如果查到就抛错
	return nil, errors.New("this car is not found")


}

type CarCtl struct {
}


func (a *CarCtl) GetCar(c *gin.Context) goft.Json {
	name := c.Query("name")

	res, err := getCar(name)
	if err != nil {
		fmt.Println("get err!")
		return gin.H{"code": "400", "error": err}
	}


	return res
}

func (a *CarCtl) CreateCar(c *gin.Context) goft.Json {
	var r *appsv1.Car
	if err := c.ShouldBindJSON(&r); err != nil {
		fmt.Println("bind json err!")
		return gin.H{"code": "400", "error": err}
	}
	res, err := createOrUpdateCar(r)
	if err != nil {
		fmt.Println("create err!")
		return gin.H{"code": "400", "error": err}
	}

	return res

}

func (a *CarCtl) UpdateCar(c *gin.Context) goft.Json {
	var r *appsv1.Car
	if err := c.ShouldBindJSON(&r); err != nil {
		fmt.Println("bind json err!")
		return gin.H{"code": "400", "error": err}
	}
	res, err := createOrUpdateCar(r)
	if err != nil {
		fmt.Println("update err!")
		return gin.H{"code": "400", "error": err}
	}

	return res

}

func (a *CarCtl) DeleteCar(c *gin.Context) goft.Json {
	name := c.Query("name")

	err := deleteCar(name)
	if err != nil {
		fmt.Println("get err!")
		return gin.H{"code": "400", "error": err}
	}


	return nil
}

func (a *CarCtl) ListCar(c *gin.Context) goft.Json {

	res, err := listCar()
	if err != nil {
		fmt.Println("list err!")
		return gin.H{"code": "400", "error": err}
	}
	return res

}

func (a *CarCtl) Name() string {
	return "CarCtl"
}

// 路由
func (a *CarCtl) Build(goft *goft.Goft) {
	// GET  http://localhost:8080/v1/car
	// GET  http://localhost:8080/v1/carlist
	// POST  http://localhost:8080/v1/car
	// DELETE  http://localhost:8080/v1/car
	// PUT  http://localhost:8080/v1/car
	goft.Handle("GET", "/v1/car", a.GetCar)
	goft.Handle("GET", "/v1/carlist", a.ListCar)
	goft.Handle("POST", "/v1/car", a.CreateCar)
	goft.Handle("DELETE", "/v1/car", a.DeleteCar)
	goft.Handle("PUT", "/v1/car", a.UpdateCar)

}

