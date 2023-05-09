package app

import (
	"flag"
	"fmt"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
	"os"
	"strings"
)

// Options 模拟配置文件
type Options struct {
	HealthPort    int
	Port          int
	NumWorker     int
	QueueCapacity int
}

func NewOptions() *Options {
	return &Options{}
}

const (
	DefaultPort       = 8080
	DefaultHealthPort = 9999
)

func (s *Options) addKlogFlags(flags *pflag.FlagSet) {
	klogFlags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	klog.InitFlags(klogFlags)

	klogFlags.VisitAll(func(f *flag.Flag) {
		f.Name = fmt.Sprintf("klog-%s", strings.ReplaceAll(f.Name, "_", "-"))
	})
	flags.AddGoFlagSet(klogFlags)
}

// AddFlags 加入命令行参数
func (s *Options) AddFlags(flags *pflag.FlagSet) {
	flags.IntVar(&s.HealthPort, "healthPort", DefaultHealthPort, "xxx")
	flags.IntVar(&s.Port, "port", DefaultPort, "xxx")

	s.addKlogFlags(flags)
}
