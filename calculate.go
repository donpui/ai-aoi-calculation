package main

func main() {

	videoHeatmap := "./examples/SALI_heatmap_example.mp4"
	videoSaliency := "./examples/SALI_saliency_example.mp4"
	outputVideo := "./output/output.mp4"
	aoiCoords := aoi_c{x_start: 400, y_start: 200, x_end: 700, y_end: 400}

	calculateVideo(videoHeatmap, videoSaliency, outputVideo, aoiCoords)

}
