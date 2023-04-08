package main

import (
	"fmt"
	"log"
	metav1 "practice_ctl/pkg/apis/meta"
	"time"

	appsv1 "practice_ctl/pkg/apis/apps/v1"
	"practice_ctl/pkg/util/stores"
	"practice_ctl/pkg/util/stores/rest"
)

func main() {
	// 配置文件
	config := &rest.Config{
		Host:    fmt.Sprintf("http://localhost:8888"),
		Timeout: time.Second,
		Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2ODA5MjE0MzAsInVzZXJuYW1lIjoiYWRtaW4ifQ.P1KJOfyF7ue0sW-Bw2-CteyrH7ro1wUjPnleV-Gnw5s",
	}
	clientSet := stores.NewForConfig(config)

	// 创建操作
	a := &appsv1.Car{
		TypeMeta: metav1.TypeMeta{
			ApiVersion: "apps/v1",
			Kind:       "Car",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "car1",
		},
		Spec: appsv1.CarSpec{
			Color: "car1",
			Brand: "car1",
			Price: "car1",
		},
	}

	c, err := clientSet.AppsV1().Car().Create(a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("name:", c.Name, "brand:", c.Spec.Brand, "color:", c.Spec.Color, "price:", c.Spec.Price)

	car1, err := clientSet.AppsV1().Car().Get("car1")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("get name: ", car1.Name)

	aaa := &appsv1.Car{
		TypeMeta: metav1.TypeMeta{
			ApiVersion: "apps/v1",
			Kind:       "Car",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "car1",
		},
		Spec: appsv1.CarSpec{
			Color: "car12222",
			Brand: "car12233",
			Price: "car12347934",
		},
	}

	carUpdate, err := clientSet.AppsV1().Car().Update(aaa)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("name: ", carUpdate.Name, "color: ", carUpdate.Spec.Color, "brand: ", carUpdate.Spec.Brand, "price: ", carUpdate.Spec.Price)

	carList, err := clientSet.AppsV1().Car().List()
	if err != nil {
		log.Fatalln(err)
	}
	for _, car := range carList.Item {
		fmt.Println(car.Name)
	}

	err = clientSet.AppsV1().Car().Delete("car1")
	if err != nil {
		log.Fatalln(err)
	}

}
