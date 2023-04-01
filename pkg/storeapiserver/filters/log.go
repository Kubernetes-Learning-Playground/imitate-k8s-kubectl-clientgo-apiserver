package filters

import (
	"log"
	"net/http"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		log.Printf("method: %s, url: %s\n", request.Method, request.URL.Path)
		log.Printf("requser host: %s", request.RemoteAddr)
		handler.ServeHTTP(response, request)
	})
}
