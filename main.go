package main

import (
	"github.com/Smotherer007/pMapTilesCutterGo/pMapTilesCutterGo"
)

func main() {
	const tileSize int = 256
	const targetPath string = "./"
	const sourcePath string = "./map.png"

	pMapTilesCutterGo.CutMapIntoTiles(tileSize, targetPath, sourcePath)
}
