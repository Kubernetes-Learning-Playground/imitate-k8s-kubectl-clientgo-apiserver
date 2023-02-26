package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Result 错误
type Result struct {
	rsp *http.Response
	err error
}

// Into 把resp转出来
func (r Result) Into(v interface{})(err error) {
	if r.err != nil {
		return r.err
	}

	body, err := ioutil.ReadAll(r.rsp.Body)
	err = json.Unmarshal(body, v)
	return err
}

