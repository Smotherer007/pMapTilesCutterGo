package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/disintegration/imaging"
)

const tileSize int = 256
const targetPath string = "/Users/patrickweppelmann/Downloads/gomap/"
const sourcePath string = "/Users/patrickweppelmann/Downloads/map.png"

func main() {

	sourceImage := loadImage(sourcePath)
	sourceImageWidth := float64(sourceImage.Bounds().Dx())
	sourceImageHeight := float64(sourceImage.Bounds().Dy())

	maxTileDim := math.Ceil(math.Max(sourceImageWidth, sourceImageHeight) / float64(tileSize))

	minZoomLevel := float64(0)
	maxZoomLevel := float64(0)
	numTilesTotalForAllLevels := float64(1)
	for math.Pow(2, maxZoomLevel) < maxTileDim {
		maxZoomLevel++
		numTilesTotalForAllLevels += math.Pow(2, 2*maxZoomLevel)
	}

	zoom := minZoomLevel
	scale := maxZoomLevel
	for zoom <= maxZoomLevel {
		canvasWidth := float64(tileSize) * math.Pow(2, zoom)
		canvasHeight := float64(tileSize) * math.Pow(2, zoom)
		canvas := createCanvas(int(canvasWidth), int(canvasHeight), color.Black)

		imageResizeWidth := scaleDimension(int(sourceImageWidth), int(scale))
		imageResizeHeight := scaleDimension(int(sourceImageHeight), int(scale))
		resizedImage := imaging.Resize(sourceImage, int(imageResizeWidth), int(imageResizeHeight), imaging.Linear)

		top := (canvasWidth - imageResizeWidth) / 2
		left := (canvasHeight - imageResizeHeight) / 2
		mergedImage := imaging.Paste(canvas, resizedImage, image.Pt(int(top), int(left)))

		createTiles(mergedImage, int(canvasWidth), int(canvasHeight), int(zoom))

		zoom++
		scale--
	}
}

func createTiles(mergedImage image.Image, width int, height int, zoom int) {
	numberOfXTiles := width / tileSize
	numberOfYTiles := height / tileSize

	for y := 0; y < numberOfYTiles; y++ {
		tileY := y * tileSize
		for x := 0; x < numberOfXTiles; x++ {
			tileX := x * tileSize
			directoryPath := fmt.Sprint(targetPath, zoom, "/", x, "/")
			if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
				err := os.MkdirAll(directoryPath, 0755)
				if err != nil {
					log.Fatal(err)
				}
			}
			upLeft := image.Pt(tileX, tileY)
			lowRight := image.Pt(tileX+tileSize, tileY+tileSize)
			tile := imaging.Crop(mergedImage, image.Rectangle{upLeft, lowRight})
			saveImage(tile, fmt.Sprint(directoryPath, y, ".png"))
		}
	}
}

func scaleDimension(dimension int, scale int) float64 {
	scaledDimension := dimension
	for i := 0; i < scale; i++ {
		scaledDimension = scaledDimension / 2
	}
	return float64(scaledDimension)
}

func createCanvas(width int, height int, color color.Color) image.Image {
	canvas := image.NewNRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			canvas.Set(x, y, color)
		}
	}
	return canvas
}

func saveImage(image image.Image, imagePath string) {
	err := imaging.Save(image, imagePath, imaging.PNGCompressionLevel(png.BestCompression))
	if err != nil {
		log.Fatal(err)
	}
}

func loadImage(imagePath string) image.Image {
	image, err := imaging.Open(imagePath, imaging.AutoOrientation(true))
	if err != nil {
		log.Fatal(err)
	}
	return image
}
