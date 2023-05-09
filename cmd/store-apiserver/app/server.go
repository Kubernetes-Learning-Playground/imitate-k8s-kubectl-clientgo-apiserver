package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
	"os"
)

func NewAPIServerCommand() *cobra.Command {
	opts := NewServerRunOptions()
	cmd := &cobra.Command{
		Use:   "store-apiserver",
		Short: "run http server",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// set default options
			completedOptions, err := Complete(opts)
			if err != nil {
				return err
			}

			// validate options
			if errs := completedOptions.Validate(); len(errs) != 0 {

				return errs[0]
			}
			stopC := make(chan struct{})
			return Run(completedOptions, stopC)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	flags := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	opts.AddFlags(flags)
	flags.Parse(os.Args[1:])
	flags.VisitAll(func(f *pflag.Flag) {
		klog.Infof("Flag: %v=%v\n", f.Name, f.Value.String())
	})

	fs := cmd.Flags()
	fs.AddFlagSet(flags)

	return cmd
}
