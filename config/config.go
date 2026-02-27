package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	PoolSize    int
	WorkerCount int
	Port        int
	PostgresDSN string
	StdOutLog   bool
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := cfg.FlagLoad(); err != nil {
		return nil, err
	}
	if err := cfg.EnvLoad(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *Config) FlagLoad() error {
	flag.IntVar(&cfg.PoolSize, "pool-size", 10, "max number of queued tasks")
	flag.IntVar(&cfg.WorkerCount, "workers", 5, "number of workers")
	flag.IntVar(&cfg.Port, "port", 8080, "http server port")
	flag.BoolVar(&cfg.StdOutLog, "stdout-log", true, "log to stdout")
	flag.StringVar(&cfg.PostgresDSN, "postgres", "postgres://database:your_password@localhost:5432/postgresdocker?sslmode=disable", "Postgres connection string")
	flag.Parse()
	return nil
}

func (cfg *Config) EnvLoad() error {
	postgresDSN := os.Getenv("POSTGRES_DSN")
	if postgresDSN != "" {
		cfg.PostgresDSN = postgresDSN
		fmt.Println("loaded postgres DSN from env")
	}
	return nil
}
