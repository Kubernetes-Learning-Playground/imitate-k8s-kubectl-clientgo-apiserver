package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
	"log"
	"net/http"
	"os"
	"practice_ctl/pkg/scheduler"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"time"
)

func NewSchedulerCommand() *cobra.Command {
	opts := NewOptions()


	cmd := &cobra.Command{
		Use:   "store-scheduler",
		Short: "store-scheduler instance",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 重要方法 runCommand
			return runCommand(opts)
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
		log.Printf("Flag: %v=%v\n", f.Name, f.Value.String())
	})

	fs := cmd.Flags()
	fs.AddFlagSet(flags)
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

	// 心跳检测健康机制
	go func() {
		handler := &healthz.Handler{
			Checks: map[string]healthz.Checker{
				"healthz": healthz.Ping,
			},
		}
		if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", opts.HealthPort), handler); err != nil {
			klog.Fatalf("Failed to start healthz endpoint: %v", err)
		}
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

const (
	DefaultNumWorker     = 3
	DefaultQueueCapacity = 10
)

func (s *Options) Config() (*Config, error) {

	c := &Config{}
	c.healthPort = s.HealthPort
	c.port = s.Port
	c.numWorker = s.NumWorker
	c.queueCapacity = s.QueueCapacity

	if s.QueueCapacity == 0 {
		s.QueueCapacity = DefaultQueueCapacity
	}

	if s.NumWorker == 0 {
		s.NumWorker = DefaultNumWorker
	}

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

	defer func() {
		sched.Stop()
		time.Sleep(time.Second * 2)
	}()

	return fmt.Errorf("finished without leader elect")
}