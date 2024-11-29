package integration

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestIntegrationBasic(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{"Image Exists in Cache", "http://localhost:8080/300/300/nginx/test_image_1.jpg",
			http.StatusOK},
		{"Image Not Found (404)", "http://localhost:8080/300/300/nginx/test_image_404.jpg",
			http.StatusNotFound},
		{"File is Not an Image", "http://localhost:8080/300/300/nginx/malicious_file.txt",
			http.StatusUnsupportedMediaType},
		{"Server Returns Error", "http://localhost:8080/300/300/nginx:8081/test_image_1.jpg",
			http.StatusInternalServerError},
		{"Image Returned", "http://localhost:8080/300/300/nginx/test_image_1.jpg", http.StatusOK},
		{"Image Smaller Than Required Size", "http://localhost:8080/1000/1000/nginx/test_image_1.jpg",
			http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{
				Timeout: time.Second * 10,
			}
			resp, err := client.Get(tt.url)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestIntegrationCache(t *testing.T) {

	imgURL := "nginx/test_image_1.jpg"
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Get("http://localhost:8080/clear")
	if err != nil {
		t.Fatalf("failed to clear cache: %v", err)
	}
	defer resp.Body.Close()

	resp, err = client.Get(fmt.Sprintf("http://localhost:8080/200/300/%s", imgURL))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, resp.Header.Get("Origin"), "http://"+imgURL)

	resp, err = client.Get(fmt.Sprintf("http://localhost:8080/200/300/%s", imgURL))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, resp.Header.Get("Origin"), "")
}
