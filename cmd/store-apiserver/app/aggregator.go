package app

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// AggregationMap 聚合API集合
// key 代表是path   value 代表是 host   --- 譬如 http://localhost:8081
type AggregationMap map[string]string

type AggregationApiServer struct {
	AggregationMap
}

func ApiProxy(targetHost string) *httputil.ReverseProxy {
	url, _ := url.Parse(targetHost)
	proxy := httputil.NewSingleHostReverseProxy(url)
	return proxy
}

func NewAggregationApiServer() *AggregationApiServer {
	aggregationMap := make(map[string]string)
	return &AggregationApiServer{AggregationMap: aggregationMap}
}

// SearchHandler 从AggregationApiServer的map中查找是否有匹配的路由
func (aa AggregationApiServer) SearchHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if server, ok := aa.AggregationMap[request.RequestURI]; ok {
			// 使用反向代理的方式去请求 外部的api server
			ApiProxy(server).ServeHTTP(response, request)
		} else {
			handler.ServeHTTP(response, request)
		}
	})
}

// SearchHandler 从AggregationApiServer的map中查找是否有匹配的路由
// 这是gin框架使用的方式，已经废弃。
//func (aa AggregationApiServer) SearchHandler(c *gin.Context) {
//	if server, ok := aa.AggregationMap[c.Request.RequestURI]; ok {
//		// 使用反向代理的方式去请求 外部的api server
//		ApiProxy(server).ServeHTTP(c.Writer, c.Request)
//	} else {
//		c.Next() //-- 中间件
//	}
//}

//
