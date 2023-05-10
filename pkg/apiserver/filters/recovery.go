package filters

import (
	"k8s.io/klog/v2"
	"net/http"
)

// RecoveryMiddleware 捕获panic中间件
func RecoveryMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				klog.Error(err)
				response.WriteHeader(http.StatusBadRequest)
				response.Write([]byte("your request make server panic, please try other."))
				return
			}
		}()
		handler.ServeHTTP(response, request)
	})
}
