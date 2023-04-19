package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <heatmap-file> <saliency-file>")
		return
	}

	videoHeatmap := os.Args[1]
	videoSaliency := os.Args[2]
	outputVideo := "./output/output.mp4"
	//aoiCoords := aoi_c{x_start: 400, y_start: 200, x_end: 700, y_end: 400}

	aoiCoords := draw(videoHeatmap)
	calculateVideo(videoHeatmap, videoSaliency, outputVideo, *aoiCoords)

}
