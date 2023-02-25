package controllers

import (
	"github.com/gorilla/websocket"
	"net/http"
)

// websocket 升级请求
var Upgrader websocket.Upgrader

func init() {
	Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}
