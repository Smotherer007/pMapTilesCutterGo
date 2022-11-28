package pMapTilesCutterGo

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"sync"

	"github.com/disintegration/imaging"
)

func CutMapIntoTiles(tileSize int, targetPath string, sourcePath string) {
	sourceImage, sourceImageWidth, sourceImageHeight := loadImage(sourcePath)

	minZoomLevel, maxZoomLevel, _ := calculateScaleParamters(sourceImage, tileSize)

	currentZoomLevel := minZoomLevel
	currentScale := maxZoomLevel

	var tilesWaitGroup sync.WaitGroup
	for currentZoomLevel <= maxZoomLevel {
		canvas, canvasWidth, canvasHeight := createCanvas(currentZoomLevel, tileSize)
		resizedImage, resizedImageWidth, resizedImageHeight := resizeImage(sourceImage, sourceImageWidth, sourceImageHeight, currentScale)
		mergedImage := mergeImageToCanvas(canvas, canvasWidth, canvasHeight, resizedImage, resizedImageWidth, resizedImageHeight)
		tilesWaitGroup.Add(1)
		go func(tileSize int, targetPath string, mergedImage image.Image, canvasWidth int, canvasHeight int, currentZoomLevel int) {
			createTiles(tileSize, targetPath, mergedImage, canvasWidth, canvasHeight, currentZoomLevel)
			tilesWaitGroup.Done()
		}(tileSize, targetPath, mergedImage, int(canvasWidth), int(canvasHeight), int(currentZoomLevel))
		currentZoomLevel++
		currentScale--
	}
	tilesWaitGroup.Wait()
}

func mergeImageToCanvas(canvas image.Image, canvasWidth int, canvasHeight int, imageToMerge image.Image, imageWidth int, imageHeight int) image.Image {
	top := (canvasWidth - imageWidth) / 2
	left := (canvasHeight - imageHeight) / 2
	mergedImage := imaging.Paste(canvas, imageToMerge, image.Pt(int(top), int(left)))
	return mergedImage
}

func resizeImage(image image.Image, imageWidth int, imageHeight int, currentScale int) (image.Image, int, int) {
	width := scaleDimension(imageWidth, currentScale)
	height := scaleDimension(imageHeight, currentScale)
	resizedImage := imaging.Resize(image, width, height, imaging.Linear)
	return resizedImage, width, height
}

func createTiles(tileSize int, targetPath string, mergedImage image.Image, width int, height int, currentZoomLevel int) {
	numberOfXTiles := width / tileSize
	numberOfYTiles := height / tileSize
	for y := 0; y < numberOfYTiles; y++ {
		tileY := y * tileSize
		for x := 0; x < numberOfXTiles; x++ {
			tileX := x * tileSize
			tile := createTile(mergedImage, tileX, tileY, tileSize)
			filePath := buildFilePath(targetPath, currentZoomLevel, x, y)
			saveImage(tile, filePath)
		}
	}
}

func createTile(mergedImage image.Image, tileX int, tileY int, tileSize int) image.Image {
	upLeft := image.Pt(tileX, tileY)
	lowRight := image.Pt(tileX+tileSize, tileY+tileSize)
	return imaging.Crop(mergedImage, image.Rectangle{upLeft, lowRight})
}

func buildFilePath(targetPath string, currentZoomLevel int, xCoordinate int, yCoordinate int) string {
	directoryPath := fmt.Sprint(targetPath, currentZoomLevel, "/", xCoordinate, "/")
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		err := os.MkdirAll(directoryPath, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	directoryPath = fmt.Sprint(directoryPath, yCoordinate, ".png")
	return directoryPath
}

func scaleDimension(dimension int, scale int) int {
	scaledDimension := dimension
	for i := 0; i < scale; i++ {
		scaledDimension = scaledDimension / 2
	}
	return scaledDimension
}

func createCanvas(currentZoomLevel int, tileSize int) (image.Image, int, int) {
	width := int(float64(tileSize) * math.Pow(2, float64(currentZoomLevel)))
	height := int(float64(tileSize) * math.Pow(2, float64(currentZoomLevel)))
	canvas := image.NewNRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			canvas.Set(x, y, color.Black)
		}
	}
	return canvas, width, height
}

func saveImage(image image.Image, imagePath string) {
	err := imaging.Save(image, imagePath, imaging.PNGCompressionLevel(png.BestCompression))
	if err != nil {
		log.Fatal(err)
	}
}

func loadImage(imagePath string) (image.Image, int, int) {
	image, err := imaging.Open(imagePath, imaging.AutoOrientation(true))
	if err != nil {
		log.Fatal(err)
	}
	width := float64(image.Bounds().Dx())
	height := float64(image.Bounds().Dy())
	return image, int(width), int(height)
}

func calculateScaleParamters(sourceImage image.Image, tileSize int) (int, int, int) {
	sourceImageWidth := float64(sourceImage.Bounds().Dx())
	sourceImageHeight := float64(sourceImage.Bounds().Dy())

	maxTileDim := math.Ceil(math.Max(sourceImageWidth, sourceImageHeight) / float64(tileSize))

	minZoomLevel := float64(0)
	maxZoomLevel := float64(0)
	numberOfTiles := float64(1)
	for math.Pow(2, maxZoomLevel) < maxTileDim {
		maxZoomLevel++
		numberOfTiles += math.Pow(2, 2*maxZoomLevel)
	}
	return int(minZoomLevel), int(maxZoomLevel), int(numberOfTiles)
}
