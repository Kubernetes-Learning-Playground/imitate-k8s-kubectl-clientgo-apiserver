package app


type Config struct {

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