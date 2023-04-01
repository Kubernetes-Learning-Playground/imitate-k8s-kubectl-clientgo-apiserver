package limit

import (
	"net/url"
	"strings"
	"sync"
)

type LimiterCache struct {
	Data sync.Map
}

var IpCache *LimiterCache


func init() {
	IpCache = &LimiterCache{}
}

// checkParam 处理query
func CheckParam(values url.Values, params string) (string, bool) {
	isParam := false
	key := ""
	// 如果没有设置params
	if params == "" {
		return key, isParam
	}
	sList := strings.Split(params, ",")
	for _, param := range sList {
		if values.Get(param) != "" {
			isParam = true
			break
		}
	}

	return key, isParam
}



