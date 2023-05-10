package filters

import (
	"k8s.io/klog/v2"
	"net/http"
	"practice_ctl/pkg/apiserver/auth"
)

func AuthorizeMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		// 登入操作，跳过授权中间件
		if request.URL.Path == "/login" {
			handler.ServeHTTP(response, request)
			return
		}
		// watch操作，跳过授权中间件。FIXME 临时方案
		if request.URL.Path == "/v1/apple/watch" || request.URL.Path == "/v1/car/watch" {
			handler.ServeHTTP(response, request)
			return
		}
		e := auth.Enforcer

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
			klog.Info("Permission verification passed")
			handler.ServeHTTP(response, request) // 进行下一步操作
		} else {
			klog.Error("Permission verification failed")
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte("the Authorize is failed"))
			return
		}

	})
}
