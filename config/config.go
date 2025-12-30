package config

type Config struct {
	PoolSize    int
	WorkerCount int
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) LoadDefaults() {
	c.PoolSize = 10
	c.WorkerCount = 3
}
