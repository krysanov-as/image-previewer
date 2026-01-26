package proxy

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHandlerProxy(t *testing.T) {
	imgPath := filepath.Join("testdata", "gopher_50x50.jpg")
	imgData, err := os.ReadFile(imgPath)
	if err != nil {
		t.Fatalf("Failed to read test image: %v", err)
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Test", "ok")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(imgData)
	}))
	defer testServer.Close()

	proxyPath := "/fill/100/100/" + testServer.URL
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", proxyPath, nil)
	Handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rr.Code)
	}

	if rr.Header().Get("X-Test") != "ok" {
		t.Errorf("Header not proxied")
	}
}
