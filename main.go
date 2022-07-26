package main

import (
	"flag"
	"fmt"

	"github.com/Smotherer007/pMapTilesCutterGo/mapTilesCutter"
)

func main() {
	flag.Usage = func() {
		fmt.Println("This tool pMapTilesCutter calculates available zoom levels and cuts an image into leaflet or google maps compatible tiles.")
		flag.PrintDefaults()
	}
	sourcePath := flag.String("sourcePath", "./map.png", "# Path of the source image / picture.")
	targetPath := flag.String("targetPath", "./", "# Path of the target folder where the tiles sould be saved.")
	tileSize := flag.Int("tileSize", 256, "# Size of the Tiles")
	aspectRatioBarsColor := flag.String("aspectRatioBarsColor", "#000000FF", "# This will change the color of the aspect ratio bars. Default value #000000FF")
	flag.Parse()
	mapTilesCutter.CutMapIntoTiles(*sourcePath, *targetPath, *tileSize, *aspectRatioBarsColor)
}
