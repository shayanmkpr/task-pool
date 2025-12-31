package config

import "flag"

type Config struct {
	PoolSize    int
	WorkerCount int
	Port        int
}

func Load() *Config {
	cfg := &Config{}
	flag.IntVar(&cfg.PoolSize, "pool-size", 10, "max number of queued tasks")
	flag.IntVar(&cfg.WorkerCount, "workers", 5, "number of workers")
	flag.IntVar(&cfg.Port, "port", 8080, "http server port")
	flag.Parse()
	return cfg
}
