package app

type Config struct {
	healthPort    int
	port          int
	numWorker     int
	queueCapacity int
}

type completedConfig struct {
	*Config
}

// CompletedConfig same as Config, just to swap private object.
type CompletedConfig struct {
	*completedConfig
}

func (c *Config) Complete() CompletedConfig {
	cc := completedConfig{c}

	return CompletedConfig{&cc}
}