// Package proxy - функции для проксирования HTTP-запросов.
package proxy

import (
	"net/http"
)

// Request выполняет HTTP-запрос к переданному URL, проксируя заголовки.
func Request(origReq *http.Request, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(origReq.Context(), origReq.Method, url, origReq.Body)
	if err != nil {
		return nil, err
	}

	for k, v := range origReq.Header {
		req.Header[k] = v
	}

	client := &http.Client{}
	return client.Do(req)
}
