package proxy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcessImage_JPEG(t *testing.T) {
	imgPath := filepath.Join("testdata", "gopher_50x50.jpg")
	data, err := os.ReadFile(imgPath)
	if err != nil {
		t.Fatalf("failed to read test image: %v", err)
	}

	resized, err := ProcessImage(data, 20, 20)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}

	if len(resized) == 0 {
		t.Errorf("Processed image is empty")
	}
}

func TestProcessImage_WrongFormat(t *testing.T) {
	data := []byte("not an image")
	_, err := ProcessImage(data, 10, 10)
	if err == nil {
		t.Errorf("Expected error for wrong format")
	}
}
