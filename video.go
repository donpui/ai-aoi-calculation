package main

import (
	"image"
	"image/color"
	"log"
	"strconv"

	gocv "gocv.io/x/gocv"
)

// Initialize the font for writing the number of white pixels
const FONT = gocv.FontHersheyPlain
const FONT_SCALE = 2.0

var FONT_COLOR = color.RGBA{255, 0, 0, 0}

var AOI_BORDER = color.RGBA{255, 0, 0, 0}

const THICKNESS = 2

type aoi_c struct {
	x_start int
	y_start int
	x_end   int
	y_end   int
}

func captureAndWriterVideo(file string, output string) (*gocv.VideoCapture, *gocv.VideoWriter, error) {
	video, err := gocv.VideoCaptureFile(file)
	if err != nil {
		log.Fatal("Error opening video file:", err)
		return nil, nil, err
	}
	defer video.Close()

	// Get the video frame rate and dimensions
	fps := int(video.Get(gocv.VideoCaptureFPS))
	width := int(video.Get(gocv.VideoCaptureFrameWidth))
	height := int(video.Get(gocv.VideoCaptureFrameHeight))

	if output != "" {
		// Create a new video writer to save the output video
		writer, err := gocv.VideoWriterFile(output, "mp4v", float64(fps), width, height, true)
		if err != nil {
			log.Fatal("Error creating video writer:", err)
			return nil, nil, err
		}
		defer writer.Close()
		return video, writer, nil
	}

	return video, nil, nil

}

func calculateVideo(heatmap string, saliency string, output string, aoiCoords aoi_c) error {
	// Open the video file

	videoS, err := gocv.VideoCaptureFile(saliency)
	if err != nil {
		log.Fatal("Error opening video file:", err)
	}
	defer videoS.Close()

	videoH, err := gocv.VideoCaptureFile(heatmap)
	if err != nil {
		log.Fatal("Error opening video file:", err)
	}
	defer videoH.Close()
	// Get the video frame rate and dimensions
	fps := int(videoH.Get(gocv.VideoCaptureFPS))
	width := int(videoH.Get(gocv.VideoCaptureFrameWidth))
	height := int(videoH.Get(gocv.VideoCaptureFrameHeight))

	// Create a new video writer to save the output video
	writer, err := gocv.VideoWriterFile(output, "mp4v", float64(fps), width, height, true)
	if err != nil {
		log.Fatal("Error creating video writer:", err)

	}
	defer writer.Close()

	// Set up a timer to write data to the frame every second
	//ticker := time.NewTicker(40 * time.Millisecond)
	//defer ticker.Stop()

	// Initialize the frame counter and white pixel count
	//frameCount := 0

	// Loop over the video frames
	for {
		// Read the next frame from the video
		frameS := gocv.NewMat()
		videoS.Read(&frameS)

		defer frameS.Close()

		frameH := gocv.NewMat()
		ok := videoH.Read(&frameH)
		if !ok {
			break
		}
		defer frameH.Close()

		// Define the area of interest
		roi := frameS.Region(image.Rect(aoiCoords.x_start, aoiCoords.y_start, aoiCoords.x_end, aoiCoords.y_end))

		// Convert the ROI to grayscale
		grayROI := gocv.NewMat()
		defer grayROI.Close()
		gocv.CvtColor(roi, &grayROI, gocv.ColorBGRToGray)

		grayROIAll := gocv.NewMat()
		defer grayROIAll.Close()
		gocv.CvtColor(frameS, &grayROIAll, gocv.ColorBGRToGray)

		// Threshold the grayscale ROI to create a binary image
		binaryROI := gocv.NewMat()
		defer binaryROI.Close()
		gocv.Threshold(grayROI, &binaryROI, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)

		binaryROIAll := gocv.NewMat()
		defer binaryROIAll.Close()
		gocv.Threshold(grayROIAll, &binaryROIAll, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)

		// Count the number of white pixels in the binary image
		numWhitePixels := gocv.CountNonZero(binaryROI)

		numWhitePixelsAll := gocv.CountNonZero(binaryROIAll)
		// Draw a rectangle around the ROI for visualization purposes
		rect := image.Rect(aoiCoords.x_end, aoiCoords.y_end, aoiCoords.x_start, aoiCoords.y_start)
		gocv.Rectangle(&frameH, rect, AOI_BORDER, 2)
		// frameCount++
		// Write the number of white pixels to the frame every second
		//select {
		//case <-ticker.C:
		// Calculate the average number of white pixels per frame
		//if frameCount > 0 {

		// Calculate percentage of white pixels in area compared to whole frame
		avgWhitePixels := float64(numWhitePixels) / float64(numWhitePixelsAll) * 100

		// Convert the number of white pixels to a string and write it to the frame
		text := strconv.FormatFloat(avgWhitePixels, 'f', 1, 64) + "%"
		gocv.PutText(&frameH, text, image.Pt(10, 50), FONT, FONT_SCALE, FONT_COLOR, THICKNESS)

		// Write the output frame to the video file
		writer.Write(frameH)
		// } else {
		// 	fmt.Println("No frames to process. Exiting...")
		// 	writer.Write(frameH)
		// }
		// Reset the frame counter and white pixel count
		// frameCount = 0
		numWhitePixels = 0
		numWhitePixelsAll = 0
		//default:
		// Do nothing if it's not time to print data to the console
		//writer.Write(frameH)
		//}
	}
	return nil
}
