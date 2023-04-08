package filters

import (
	"net/http"
	"practice_ctl/pkg/storeapiserver/auth"
)

// AuthenticateMiddleware 认证中间件
func AuthenticateMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		// 如果是登入接口，跳过认证中间件
		if request.URL.Path == "/login" {
			handler.ServeHTTP(response, request)
			return
		}
		// 认证
		err := auth.ValidateToken(response, request)
		if err == nil {
			handler.ServeHTTP(response, request)
		} else {
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte("the authenticate is failed"))
		}

	})
}
