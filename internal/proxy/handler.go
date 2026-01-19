// Package proxy - функции для проксирования HTTP-запросов.
package proxy

import (
	"io"
	"log"
	"net/http"
	"strings"
)

// Handler обрабатывает запросы /fill/{width}/{height}/{url}.
func Handler(w http.ResponseWriter, r *http.Request) {
	previewParams := strings.SplitN(r.URL.Path, "/", 5)
	if len(previewParams) < 5 {
		http.Error(w, "invalid URL format", http.StatusBadRequest)
		return
	}

	width := previewParams[2]
	height := previewParams[3]
	rawURL := previewParams[4]

	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	log.Printf("Requested preview: %sx%s for %s", width, height, rawURL)

	resp, err := Request(r, rawURL)
	if err != nil {
		log.Printf("Proxy error: %v", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Failed to write response body: %v", err)
	}
}
