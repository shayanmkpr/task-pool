package config

type Config struct {
	PoolSize    int
	WorkerCount int
	Port        int
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) LoadDefaults() {
	c.PoolSize = 4
	c.WorkerCount = 3
	c.Port = 8080
}
