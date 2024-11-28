package image_transformer

import (
	"golang.org/x/image/draw"
	"image"
	_ "image/jpeg"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func Resize(img image.Image, width, height int) (image.Image, error) {

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	cropped := crop(img, width, height)

	draw.ApproxBiLinear.Scale(dst, dst.Rect, cropped, cropped.Bounds(), draw.Over, nil)

	return dst, nil
}

func crop(img image.Image, width, height int) image.Image {
	croppedWidth, croppedHeight := getCroppedSizes(img.Bounds().Dx(), img.Bounds().Dy(), width, height)

	x0 := (img.Bounds().Dx() - croppedWidth) / 2
	y0 := (img.Bounds().Dy() - croppedHeight) / 2

	croppedImg := image.Rect(x0, y0, x0+croppedWidth, y0+croppedHeight)

	return img.(SubImager).SubImage(croppedImg)
}

func getCroppedSizes(srcWidth, srcHeight, dstWidth, dstHeight int) (croppedWidth, croppedHeight int) {
	if croppedWidth = dstWidth * srcHeight / dstHeight; croppedWidth <= srcWidth {
		return croppedWidth, srcHeight
	}

	return srcWidth, dstHeight * srcWidth / dstWidth
}
