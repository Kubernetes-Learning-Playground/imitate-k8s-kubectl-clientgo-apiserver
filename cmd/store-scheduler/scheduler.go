package main

import (
	"k8s.io/component-base/cli"
	"os"
	"practice_ctl/cmd/store-scheduler/app"
)

func main() {
	command := app.NewSchedulerCommand()
	code := cli.Run(command)
	os.Exit(code)
}
