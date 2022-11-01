package main

import (
	"github.com/shenyisyn/goft-gin/goft"
	"practice_ctl/pkg/blogapiserver/controllers"
)

func main() {
	goft.Ignite().Config().
		Mount("",
			&controllers.VersionCtl{},
		).
		LaunchWithPort(8080)
}
