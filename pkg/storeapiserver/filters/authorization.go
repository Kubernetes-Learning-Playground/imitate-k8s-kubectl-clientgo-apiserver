package filters

import (
	"fmt"
	"github.com/casbin/casbin"
	"net/http"
)

var Enforcer *casbin.Enforcer

func AuthorizeMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		e := Enforcer

		////从DB加载策略
		//e.LoadPolicy()

		// 获取请求的URI
		obj := request.URL.RequestURI()
		// 获取请求方法
		act := request.Method
		// 获取用户的角色 从header拿取
		sub := request.Header.Get("username")

		// 判断策略中是否存在
		if ok := e.Enforce(sub, obj, act); ok {
			fmt.Println("恭喜您,权限验证通过")
			handler.ServeHTTP(response, request) // 进行下一步操作
		} else {
			fmt.Println("很遗憾,权限验证没有通过")
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte("the Authorize is failed"))
		}

	})
}