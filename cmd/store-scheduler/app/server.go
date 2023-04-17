package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"practice_ctl/pkg/scheduler"
)

func NewSchedulerCommand() *cobra.Command {
	s := NewOptions()
	cmd := &cobra.Command{
		Use:   "store-scheduler",
		Short: "store-scheduler instance",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 重要方法 runCommand
			return runCommand(s)
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

	return cmd
}

// runCommand runs the scheduler.
func runCommand(opts *Options) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		stopCh := make(chan struct{})
		<-stopCh
		cancel()
	}()

	cc, sched, err := SetUp(ctx, opts)
	if err != nil {
		return err
	}

	return Run(ctx, cc, sched)

}

func (s *Options) Validate() []error {
	var errs []error

	return errs
}

func (s *Options) Config() (*Config, error) {

	c := &Config{}

	return c, nil
}


func SetUp(ctx context.Context, opts *Options) (*CompletedConfig, *scheduler.Scheduler, error) {
	// opt验证配置
	if errs := opts.Validate(); len(errs) > 0 {
		return nil, nil, errors.New("validate error")
	}
	// 配置对象
	c, err := opts.Config()
	if err != nil {
		return nil, nil, err
	}

	// Get the completed config
	// 完成需要用的配置文件
	cc := c.Complete()

	sched, err := scheduler.New()

	return &cc, sched, nil
}

func Run(ctx context.Context, cc *CompletedConfig, sched *scheduler.Scheduler) error {

	sched.Run(ctx)

	return fmt.Errorf("finished without leader elect")
}