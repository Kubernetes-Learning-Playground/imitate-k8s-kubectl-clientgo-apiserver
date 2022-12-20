package main

import (
	"fmt"
	"log"
	v1 "practice_ctl/pkg/apis/core/v1"

	//v1 "practice_ctl/pkg/apis/core/v1"

	//"log"

	"practice_ctl/pkg/util/stores"
	"practice_ctl/pkg/util/stores/rest"
	"time"
)


func main() {
	// 配置文件
	config := &rest.Config{
		Host:    fmt.Sprintf("http://localhost:8080"),
		Timeout: time.Second,
	}
	clientSet := stores.NewForConfig(config)

	// 创建操作
	a := &v1.Apple{
		Name: "apple1",
		Size: "apple1",
		Color: "apple1",
		Place: "apple1",
		Price: "apple1",
	}
	c, err := clientSet.CoreV1().Apple().Create(a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("name:", c.Name,  "size:", c.Size, "color:", c.Color, "place:", c.Place, "price:", c.Price)

	apple1, err := clientSet.CoreV1().Apple().Get(c.Name)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("name: ", apple1.Name)

	aaa := &v1.Apple{
		Name: "apple1",
		Size: "apple1dddd",
		Color: "apple1ccc",
		Place: "apple1ccc",
		Price: "apple1ccc",
	}

	appleupdate, err := clientSet.CoreV1().Apple().Update(aaa)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("name: ", appleupdate.Name,  "size: ", appleupdate.Size, "color: ", appleupdate.Color, "place: ", appleupdate.Place, "price: ", appleupdate.Price)

	appleList, err := clientSet.CoreV1().Apple().List()
	if err != nil {
		log.Fatalln(err)
	}
	for _, apple := range appleList.Item {
		fmt.Println(apple.Name)
	}


}
