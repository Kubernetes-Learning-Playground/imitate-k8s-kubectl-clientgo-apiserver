package auth

import (
	"fmt"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

var Enforcer *casbin.Enforcer

//func init() {
//	CasbinSetup()
//}

// 初始化casbin
func CasbinSetup(){

	e := casbin.NewEnforcer("./conf/rbac_models.conf")


	Enforcer = e
	//return e
}


func Hello(c *gin.Context) {
	fmt.Println("Hello 接收到GET请求..")
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": "Hello 接收到GET请求..",
	})
}


func main() {
	//获取router路由对象
	r := gin.New()


	//增加policy
	r.GET("/api/v1/add", func(c *gin.Context) {
		fmt.Println("增加Policy")
		if ok := Enforcer.AddPolicy("admin", "/api/v1/world", "GET"); !ok {
			fmt.Println("Policy已经存在")
		} else {
			fmt.Println("增加成功")
		}
	})
	//删除policy
	r.DELETE("/api/v1/delete", func(c *gin.Context) {
		fmt.Println("删除Policy")
		if ok := Enforcer.RemovePolicy("admin", "/api/v1/world", "GET"); !ok {
			fmt.Println("Policy不存在")
		} else {
			fmt.Println("删除成功")
		}
	})
	//获取policy
	r.GET("/api/v1/get", func(c *gin.Context) {
		fmt.Println("查看policy")
		list := Enforcer.GetPolicy()
		for _, vlist := range list {
			for _, v := range vlist {
				fmt.Printf("value: %s, ", v)
			}
		}
	})

	//使用自定义拦截器中间件
	//r.Use(Authorize())




	//创建请求
	r.GET("/api/v1/hello", Hello)
	r.Run(":9000") //参数为空 默认监听8080端口
}

