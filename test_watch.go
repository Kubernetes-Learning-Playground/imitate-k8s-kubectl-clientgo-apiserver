package main

import (
	"fmt"
	"practice_ctl/pkg/util/stores"
	"practice_ctl/pkg/util/stores/rest"
	"time"
)

//
func main() {
	// 配置文件
	config := &rest.Config{
		Host:    fmt.Sprintf("http://localhost:8080"),
		Timeout: time.Second,
	}
	clientSet := stores.NewForConfig(config)

	res := clientSet.CoreV1().Apple().Watch()
	for i := range res.WChan {
		r := i.([]byte)
		fmt.Println(string(r))
	}
}
