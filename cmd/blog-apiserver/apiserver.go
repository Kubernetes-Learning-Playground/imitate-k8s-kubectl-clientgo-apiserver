package main

import (
	"github.com/shenyisyn/goft-gin/goft"
	"practice_ctl/pkg/blogapiserver/controllers"
)

//本课程来自 程序员在囧途(www.jtthink.com) 咨询群：98514334
func main() {
	goft.Ignite().Config().
		Mount("",
			&controllers.VersionCtl{},
		).
		LaunchWithPort(8080)
}
