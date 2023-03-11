package main

import "github.com/gin-gonic/gin"


// 测试 aggregator api server的扩展 server
func main() {
	r := gin.New()

	defer func() {
		r.Run(":8081")
	}()


	r.GET("/test-aggregator", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test-aggregator-success"})
	})


}