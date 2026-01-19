package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello World"))
	}))
	defer testServer.Close()

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := Request(req, testServer.URL)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Hello World" {
		t.Errorf("Unexpected body: %s", body)
	}
}

func TestHandlerBadURL(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fill/100", nil)
	Handler(responseRecorder, req)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 BadRequest, got %d", responseRecorder.Code)
	}
}

func TestHandlerProxy(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Test", "ok")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Proxy works"))
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
