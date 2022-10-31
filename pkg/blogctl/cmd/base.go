package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"practice_ctl/pkg/blogctl/config"
	"log"
)

// RunCmd 本课程来自 程序员在囧途(www.jtthink.com) 咨询群：98514334
func RunCmd() {
	cmd := &cobra.Command{
		Use:          "blogctl",
		Short:        "程序员在囧途个人博客",
		Example:      "blogctl",
		SilenceUsage: true,
	}

	cfg := config.LoadConfigFile()
	fmt.Println(cfg)
	//加入子命令
	cmd.AddCommand(versionCmd)
	err := cmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
