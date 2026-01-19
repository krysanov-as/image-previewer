// Package main - запуск HTTP-сервера для проксирования изображений.
package main

import (
	"log"
	"net/http"

	"github.com/krysanov-as/img-previewer/internal/config" //nolint:depguard
	"github.com/krysanov-as/img-previewer/internal/proxy"  //nolint:depguard
)

func main() {
	cfg := config.NewDefaultConfig()

	http.HandleFunc("/fill/", proxy.Handler)

	srv := &http.Server{
		Addr:         cfg.Port,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	}

	log.Printf("Starting Previewer server on %s with cache size %d\n", cfg.Port, cfg.CacheSize)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server start failed: %v", err)
	}
}
