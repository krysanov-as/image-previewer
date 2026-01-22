// Package config содержит конфигурацию сервиса previewer.
package config

import (
	"flag"
	"time"
)

type Config struct {
	Port      string
	Timeout   time.Duration
	CacheSize int64
	CacheDir  string
}

func GetConfig() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.Port, "port", ":8080", "HTTP server port")
	flag.DurationVar(&cfg.Timeout, "timeout", 10*time.Second, "HTTP server timeout")
	flag.Int64Var(&cfg.CacheSize, "cache", 104857600, "Cache size in bytes")
	flag.StringVar(&cfg.CacheDir, "dir", "./cache", "Cache directory")
	flag.Parse()
	return cfg
}
