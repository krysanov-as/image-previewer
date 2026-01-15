// Package config содержит конфигурацию сервиса previewer.
package config

import "time"

// Config структура конфига.
type Config struct {
	Port      string
	CacheSize int // Максимальный размер LRU-кэша (кол-во изображений)
	Timeout   time.Duration
}

// NewDefaultConfig - дефолтные настройки сервиса.
func NewDefaultConfig() *Config {
	return &Config{
		Port:      ":8080",
		CacheSize: 10,
		Timeout:   10 * time.Second,
	}
}
