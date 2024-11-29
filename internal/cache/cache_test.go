package cache

import (
	"image"
	"image/color"
	"os"
	"testing"
)

func createTestImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	var x, y uint8
	for x = 0; x < 100; x++ {
		for y = 0; y < 100; y++ {
			img.Set(int(x), int(y), color.RGBA{x, y, 255, 255})
		}
	}
	return img
}

func TestLruImageCache(t *testing.T) {
	dir, err := os.MkdirTemp("", "cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(dir)

	cache := NewLruImageCache(2, dir)

	img1 := createTestImage()
	img2 := createTestImage()
	img3 := createTestImage()

	err = cache.Save("image1.jpg", img1)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	imgRetrieved := cache.Get("image1.jpg")
	if imgRetrieved == nil {
		t.Error("Expected to retrieve image1.jpg, got nil")
	}

	err = cache.Save("image2.jpg", img2)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	imgRetrieved = cache.Get("image2.jpg")
	if imgRetrieved == nil {
		t.Error("Expected to retrieve image2.jpg, got nil")
	}

	err = cache.Save("image3.jpg", img3) // This should evict image1
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	imgRetrieved = cache.Get("image1.jpg")
	if imgRetrieved != nil {
		t.Error("Expected image1.jpg to be evicted, got non-nil")
	}

	imgRetrieved = cache.Get("image2.jpg")
	if imgRetrieved == nil {
		t.Error("Expected to retrieve image2.jpg, got nil")
	}

	imgRetrieved = cache.Get("image3.jpg")
	if imgRetrieved == nil {
		t.Error("Expected to retrieve image3.jpg, got nil")
	}
}

func TestLruImageCache_Load(t *testing.T) {
	dir, err := os.MkdirTemp("", "cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(dir)

	cache := NewLruImageCache(2, dir)

	img := createTestImage()
	err = cache.Save("image.jpg", img)
	if err != nil {
		t.Errorf("Expected no error while saving, got: %v", err)
	}

	cache2 := NewLruImageCache(2, dir)
	err = cache2.Load()
	if err != nil {
		t.Errorf("Expected no error while loading, got: %v", err)
	}

	imgRetrieved := cache2.Get("image.jpg")
	if imgRetrieved == nil {
		t.Error("Expected to retrieve image.jpg after loading, got nil")
	}
}

func TestLruImageCache_SaveError(t *testing.T) {
	dir, err := os.MkdirTemp("", "cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(dir)

	cache := NewLruImageCache(2, dir)

	img := createTestImage()

	os.RemoveAll(dir)

	err = cache.Save("image.jpg", img)
	if err == nil {
		t.Error("Expected an error when saving to a non-existent directory, got nil")
	}
}

func TestLruImageCache_LoadError(t *testing.T) {
	dir, err := os.MkdirTemp("", "cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(dir)

	cache := NewLruImageCache(2, dir)

	// Attempt loading from an empty directory
	err = cache.Load()
	if err != nil {
		t.Errorf("Expected no error when loading from empty directory, got: %v", err)
	}
}
