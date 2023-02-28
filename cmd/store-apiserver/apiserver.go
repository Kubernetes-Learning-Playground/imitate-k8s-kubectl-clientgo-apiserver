package main

import (
	"github.com/shenyisyn/goft-gin/goft"
	"practice_ctl/pkg/storeapiserver/controllers"
)

func main() {

	goft.Ignite().Config().
		Mount("",
			&controllers.VersionCtl{},
			&controllers.AppleCtl{},
			&controllers.CarCtl{},
		).LaunchWithPort(8080)
}
