// Package main - запуск HTTP-сервера для проксирования изображений.
package main

import (
	"log"
	"net/http"

	"github.com/krysanov-as/img-previewer/internal/config" //nolint:depguard
	"github.com/krysanov-as/img-previewer/internal/proxy"  //nolint:depguard
)

func main() {
	cfg := config.GetConfig()

	if err := proxy.InitCache(cfg.CacheSize, cfg.CacheDir); err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}

	http.HandleFunc("/fill/", proxy.Handler)

	srv := &http.Server{
		Addr:         cfg.Port,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	}

	log.Printf(
		"Starting Previewer server on %s, cache dir: %s, max cache: %d bytes",
		cfg.Port, cfg.CacheDir, cfg.CacheSize,
	)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server start failed: %v", err)
	}
}
