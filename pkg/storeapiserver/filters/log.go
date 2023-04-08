package filters

import (
	"log"
	"net/http"
	"time"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		startTime := time.Now()
		handler.ServeHTTP(response, request)
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		log.Printf("request time: %s", latencyTime)  // 执行时间
		log.Printf("method: %s, url: %s\n", request.Method, request.URL.Path) // 请求方法 请求url
		log.Printf("requser host: %s", request.RemoteAddr)	// 请求host
	})
}
