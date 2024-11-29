package imagetransformer

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path"
	"strconv"
	"testing"
)

type testImage struct {
	bounds image.Rectangle
}

func (ti *testImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (ti *testImage) Bounds() image.Rectangle {
	return ti.bounds
}

func (ti *testImage) At(_, _ int) color.Color {
	return color.RGBA{255, 0, 0, 255} // Красный цвет для тестирования
}

func (ti *testImage) SubImage(r image.Rectangle) image.Image {
	return &testImage{bounds: r}
}

func createTestImage(width, height int) image.Image {
	return &testImage{bounds: image.Rect(0, 0, width, height)}
}

func TestResize(t *testing.T) {
	img := createTestImage(100, 100)

	resizedImg, err := Resize(img, 50, 50)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resizedImg.Bounds().Dx() != 50 || resizedImg.Bounds().Dy() != 50 {
		t.Errorf("Expected resized image dimensions to be 50x50, got: %dx%d", resizedImg.Bounds().Dx(),
			resizedImg.Bounds().Dy())
	}
}

func TestResizeAspectRatio(t *testing.T) {
	img := createTestImage(300, 200)

	resizedImg, err := Resize(img, 150, 100)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resizedImg.Bounds().Dx() != 150 || resizedImg.Bounds().Dy() != 100 {
		t.Errorf("Expected resized image dimensions to be 150x100, got: %dx%d", resizedImg.Bounds().Dx(),
			resizedImg.Bounds().Dy())
	}
}

func TestResizeCrop(t *testing.T) {
	img := createTestImage(400, 300)

	resizedImg, err := Resize(img, 200, 100)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resizedImg.Bounds().Dx() != 200 || resizedImg.Bounds().Dy() != 100 {
		t.Errorf("Expected resized image dimensions to be 200x100, got: %dx%d", resizedImg.Bounds().Dx(),
			resizedImg.Bounds().Dy())
	}
}

func TestGetCroppedSizes(t *testing.T) {
	tests := []struct {
		srcWidth, srcHeight, dstWidth, dstHeight    int
		expectedCroppedWidth, expectedCroppedHeight int
	}{
		{400, 300, 200, 100, 400, 200},
		{400, 300, 100, 100, 300, 300},
		{400, 300, 100, 50, 400, 200},
		{400, 300, 50, 100, 150, 300},
	}

	for _, tt := range tests {
		croppedWidth, croppedHeight := getCroppedSizes(tt.srcWidth, tt.srcHeight, tt.dstWidth, tt.dstHeight)
		if croppedWidth != tt.expectedCroppedWidth || croppedHeight != tt.expectedCroppedHeight {
			t.Errorf("For src %dx%d and dst %dx%d, expected cropped size %dx%d, got %dx%d",
				tt.srcWidth, tt.srcHeight, tt.dstWidth, tt.dstHeight,
				tt.expectedCroppedWidth, tt.expectedCroppedHeight,
				croppedWidth, croppedHeight)
		}
	}
}

func TestCrop(t *testing.T) {
	img := createTestImage(100, 100)

	croppedImg := crop(img, 50, 50)
	if croppedImg.Bounds().Dx() != 100 || croppedImg.Bounds().Dy() != 100 {
		t.Errorf("Expected cropped image dimensions to be 50x50, got: %dx%d", croppedImg.Bounds().Dx(),
			croppedImg.Bounds().Dy())
	}

	if _, ok := croppedImg.(SubImager); !ok {
		t.Error("Expected cropped image to implement SubImager")
	}
}

func TestResizeGopherImage(t *testing.T) {
	type Size struct {
		Width  int
		Height int
	}

	inputFile := "testdata/_gopher_original_1024x504.jpg"
	outputDir := "output"

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	defer os.RemoveAll(outputDir)

	sizes := []Size{
		{Width: 1024, Height: 252},
		{Width: 2000, Height: 1000},
		{Width: 200, Height: 700},
		{Width: 256, Height: 126},
		{Width: 333, Height: 666},
		{Width: 500, Height: 500},
		{Width: 50, Height: 50},
	}

	imgFile, err := os.Open(inputFile)
	if err != nil {
		t.Fatalf("Failed to open image file: %v", err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		t.Fatalf("Failed to decode image: %v", err)
	}

	for _, size := range sizes {
		width := size.Width
		height := size.Height

		resizedImg, err := Resize(img, width, height)
		if err != nil {
			t.Errorf("Failed to resize image to %dx%d: %v", width, height, err)
			continue
		}

		outputFilePath := path.Join(outputDir,
			"resized_"+path.Base(inputFile[:len(inputFile)-4])+
				"_"+strconv.Itoa(width)+"x"+strconv.Itoa(height)+".jpg")
		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			t.Errorf("Failed to create output file: %v", err)
			continue
		}

		if err := jpeg.Encode(outputFile, resizedImg, nil); err != nil {
			t.Errorf("Failed to save resized image: %v", err)
		}

		if resizedImg.Bounds().Dx() != width || resizedImg.Bounds().Dy() != height {
			t.Errorf("Resizing error for %dx%d: got %dx%d", width, height, resizedImg.Bounds().Dx(),
				resizedImg.Bounds().Dy())
		}
	}
}
