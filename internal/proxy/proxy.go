package proxy

import (
	"net/http"
)

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
