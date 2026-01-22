package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Handler обрабатывает запросы /fill/{width}/{height}/{url}.
func Handler(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(strings.Trim(r.URL.Path, "/"), "/", 4)
	if len(parts) < 4 {
		http.Error(w, "invalid URL format: expected /fill/width/height/url", http.StatusBadRequest)
		return
	}

	width, err := strconv.Atoi(parts[1])
	if err != nil || width <= 0 {
		http.Error(w, "invalid width", http.StatusBadRequest)
		return
	}

	height, err := strconv.Atoi(parts[2])
	if err != nil || height <= 0 {
		http.Error(w, "invalid height", http.StatusBadRequest)
		return
	}

	rawURL := parts[3]
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	cacheKey := fmt.Sprintf("%dx%d-%s", width, height, rawURL)

	if imageCache != nil {
		if cachedData, ok := imageCache.Get(cacheKey); ok {
			log.Printf("Cache hit for: %s", cacheKey)
			w.Header().Set("Content-Type", "image/jpeg")
			w.Write(cachedData)
			return
		}
	}

	log.Printf("Processing preview: %dx%d for %s", width, height, rawURL)

	resp, err := Request(r, rawURL)
	if err != nil {
		log.Printf("Failed to fetch image: %v", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Remote server returned: %d", resp.StatusCode)
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	processedImage, err := ProcessImage(body, width, height)
	if err != nil {
		log.Printf("Image processing failed: %v", err)
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}

	if imageCache != nil {
		if err := imageCache.Set(cacheKey, processedImage); err != nil {
			log.Printf("Failed to cache image: %v", err)
		} else {
			log.Printf("Saved to cache: %s", cacheKey)
		}
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(processedImage)
}
