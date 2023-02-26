package main

import (
	"fmt"

	"log"
	v1 "practice_ctl/pkg/apis/core/v1"
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
		ApiVersion: "core/v1",
		Kind: "APPLE",
		Metadata: v1.Metadata{
			Name: "applexxxxxxx",
		},
		Spec: v1.AppleSpec{
			Size: "apple1",
			Color: "apple1",
			Place: "apple1",
			Price: "apple1",
		},
		Status: v1.AppleStatus{
			Status: "created",
		},

	}
	c, err := clientSet.CoreV1().Apple().Create(a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("name:", c.Name,  "size:", c.Spec.Size, "color:", c.Spec.Color, "place:", c.Spec.Place, "price:", c.Spec.Price)

	apple1, err := clientSet.CoreV1().Apple().Get(c.Name)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("name: ", apple1.Name)

	aaa := &v1.Apple{
		ApiVersion: "core/v1",
		Kind: "APPLE",
		Metadata: v1.Metadata{
			Name: "apple-test11",
		},
		Spec: v1.AppleSpec{
			Size: "apple1dddd",
			Color: "apple1ccc",
			Place: "apple1ccc",
			Price: "apple1ccc",
		},
	}

	appleUpdate, err := clientSet.CoreV1().Apple().Update(aaa)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("name: ", appleUpdate.Name,  "size: ", appleUpdate.Spec.Size, "color: ", appleUpdate.Spec.Color, "place: ", appleUpdate.Spec.Place, "price: ", appleUpdate.Spec.Price)

	appleList, err := clientSet.CoreV1().Apple().List()
	if err != nil {
		log.Fatalln(err)
	}
	for _, apple := range appleList.Item {
		fmt.Println(apple.Name)
	}


}
