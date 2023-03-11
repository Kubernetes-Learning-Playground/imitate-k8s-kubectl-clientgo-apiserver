package aggregator

import (
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"net/url"
)

// 聚合API集合
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
func (aa AggregationApiServer) SearchHandler(c *gin.Context) {
	if server, ok := aa.AggregationMap[c.Request.RequestURI]; ok {
		// 使用反向代理的方式去请求 外部的api server
		ApiProxy(server).ServeHTTP(c.Writer, c.Request)
	} else {
		c.Next() //-- 中间件
	}
}

//

