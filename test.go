package main

import (
	"fmt"
	"k8sblog/pkg/util/blogs"
	"k8sblog/pkg/util/blogs/rest"
	"log"
	"time"
)

// 本课程来自 程序员在囧途(www.jtthink.com) 咨询群：98514334
func main() {
	config := &rest.Config{
		Host:    "http://localhost:8080",
		Timeout: time.Second,
	}
	//v := &unverstioned.Version{}
	clientSet := blogs.NewForConfig(config)
	v, err := clientSet.Core().Version().Get()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(v)

}

// 本课程来自 程序员在囧途(www.jtthink.com) 咨询群：98514334
