package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
)

type ClientVersionInfo struct {
	Version   string
	GoVersion string
}

func (cv *ClientVersionInfo) String() string {
	return fmt.Sprintf("versionInfo:{version:%s,goversion:%s}", cv.Version, cv.GoVersion)
}

func NewClientVersionInfo() *ClientVersionInfo {
	return &ClientVersionInfo{
		Version:   "0.1",
		GoVersion: "go1.18",
	}
}

var versionCmd = &cobra.Command{
	Use:          "version", // 以  xxxx version 运行
	Short:        "v",
	Example:      "storectl version",
	SilenceUsage: true,
	RunE: func(c *cobra.Command, args []string) error {
		fmt.Printf("Client Version:%s\n", NewClientVersionInfo())

		return nil
	},
}
