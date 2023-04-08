package filters

import (
	"context"
	"net/http"
	"time"
)

// RequestTimeoutMiddleware 请求超时中间件
func RequestTimeoutMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		// FIXME: 超时配置设置成可配置的
		ctx, cancel := context.WithTimeout(request.Context(), time.Second*2)

		defer func() {
			if ctx.Err() == context.DeadlineExceeded {
				response.WriteHeader(http.StatusGatewayTimeout)
				response.Write([]byte("request timeout, please try angin"))
			}
			// clear
			cancel()
		}()

		request = request.WithContext(ctx)

		handler.ServeHTTP(response, request)
	})
}
