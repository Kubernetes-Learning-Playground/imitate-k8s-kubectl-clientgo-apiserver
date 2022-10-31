package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
)

type ClientVersionInfo struct {
	Version   string
	GoVersion string
}

//本课程来自 程序员在囧途(www.jtthink.com) 咨询群：98514334
func (cv *ClientVersionInfo) String() string {
	return fmt.Sprintf("versionInfo:{version:%s,goversion:%s}", cv.Version, cv.GoVersion)
}

//本课程来自 程序员在囧途(www.jtthink.com) 咨询群：98514334
func NewClientVersionInfo() *ClientVersionInfo {
	return &ClientVersionInfo{
		Version:   "0.1",
		GoVersion: "go1.18",
	}
}

var versionCmd = &cobra.Command{
	Use:          "version", // 这意味着 可以  xxxx version 运行
	Short:        "v",
	Example:      "blogctl version",
	SilenceUsage: true,
	RunE: func(c *cobra.Command, args []string) error {
		fmt.Printf("Client Version:%s\n", NewClientVersionInfo())

		return nil
	},
}

//本课程来自 程序员在囧途(www.jtthink.com) 咨询群：98514334
