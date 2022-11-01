package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"practice_ctl/pkg/blogctl/config"
	"log"
)

// RunCmd
func RunCmd() {
	cmd := &cobra.Command{
		Use:          "blogctl",
		Short:        "模仿kubectl",
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
