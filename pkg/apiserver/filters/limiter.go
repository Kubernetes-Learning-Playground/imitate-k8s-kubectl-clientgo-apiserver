package filters

import (
	"k8s.io/klog/v2"
	"net/http"
	"practice_ctl/pkg/apiserver/filters/limit"
)

// IpLimiterMiddleware ip限流中间件
func IpLimiterMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		ip := request.RemoteAddr
		klog.Info("ip: ", ip)

		var limiter *limit.Bucket

		if v, ok := limit.IpCache.Data.Load(ip); ok {
			limiter = v.(*limit.Bucket)
		} else {
			limiter = limit.NewBucket(limit.DefaultCap, limit.DefaultRate)
			limit.IpCache.Data.Store(ip, limiter)
		}

		// 如果限流器接受，则走到下一个中间件，不然就报错
		if limiter.IsAccept() {
			handler.ServeHTTP(response, request)
		} else {
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte("this ip is too many request!!"))
			return
		}
	})
}

// ParamLimiterMiddleware query限流中间件
//func ParamLimiterMiddleware(next http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, req *http.Request) {
//
//		if key, ok := limit.CheckParam(req.URL.Query(), sysconfig.SysConfig1.Server.Params); ok {
//
//			var limiter *limit.Bucket
//
//			if v, ok := limit.IpCache.Data.Load(key); ok {
//				limiter = v.(*limit.Bucket)
//			} else {
//				limiter = limit.NewBucket(1, limit.DefaultRate)
//				limit.IpCache.Data.Store(key, limiter)
//			}
//
//			if limiter.IsAccept() {
//				next(w, req)
//			} else {
//				w.Write([]byte("this query is too many request!!"))
//			}
//		} else {
//			next(w, req)
//		}
//
//	}
//}
