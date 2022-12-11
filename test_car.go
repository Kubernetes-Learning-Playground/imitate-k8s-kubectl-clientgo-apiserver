package main

import (
	"fmt"
	"log"
	"time"

	appsv1 "practice_ctl/pkg/apis/apps/v1"
	"practice_ctl/pkg/util/stores"
	"practice_ctl/pkg/util/stores/rest"

)


func main() {
	// 配置文件
	config := &rest.Config{
		Host:    fmt.Sprintf("http://localhost:8080"),
		Timeout: time.Second,
	}
	clientSet := stores.NewForConfig(config)

	// 创建操作
	a := &appsv1.Car{
		Name: "car1",
		Color: "car1",
		Brand: "car1",
		Price: "car1",
	}
	c, err := clientSet.Apps().Car().Create(a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("name:", c.Name,  "brand:", c.Brand, "color:", c.Color, "price:", c.Price)

	car1, err := clientSet.Apps().Car().Get(c.Name)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("name: ", car1.Name)

	aaa := &appsv1.Car{
		Name: "car1",
		Color: "car1ccc",
		Brand: "car1ccc",
		Price: "car1ccc",
	}

	carupdate, err := clientSet.Apps().Car().Update(aaa)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("name: ", carupdate.Name, "color: ", carupdate.Color, "brand: ", carupdate.Brand, "price: ", carupdate.Price)

	carList, err := clientSet.Apps().Car().List()
	if err != nil {
		log.Fatalln(err)
	}
	for _, car := range carList.Item {
		fmt.Println(car.Name)
	}


}

