package main

import (
	"fmt"
	"log"
	"practice_ctl/pkg/util/blogs"
	"practice_ctl/pkg/util/blogs/rest"
	"time"
)


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
