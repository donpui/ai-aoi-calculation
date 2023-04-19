package main

import (
	"fmt"
	"image"
	"image/color"

	gocv "gocv.io/x/gocv"
)

func imageCalculation() {
	// Load the saliency image
	img := gocv.IMRead("examples/saliency_shopify.jpg", gocv.IMReadColor)

	// Define the area of interest
	x1, y1, x2, y2 := 100, 100, 300, 300
	roi := img.Region(image.Rect(x1, y1, x2, y2))

	grayFullROI := gocv.NewMat()
	defer grayFullROI.Close()
	gocv.CvtColor(img, &grayFullROI, gocv.ColorBGRToGray)

	// Convert the ROI to grayscale
	grayROI := gocv.NewMat()
	defer grayROI.Close()
	gocv.CvtColor(roi, &grayROI, gocv.ColorBGRToGray)

	// Threshold the grayscale ROI to create a binary image
	binaryROI := gocv.NewMat()
	defer binaryROI.Close()
	gocv.Threshold(grayROI, &binaryROI, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)

	// Threshold the grayscale ROI to create a binary image
	binaryFullROI := gocv.NewMat()
	defer binaryFullROI.Close()
	gocv.Threshold(grayFullROI, &binaryFullROI, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)

	// Count the number of white pixels in the binary image
	numWhitePixels := gocv.CountNonZero(binaryROI)
	numWhitePixelsAll := gocv.CountNonZero(binaryFullROI)

	// Draw a rectangle around the ROI for visualization purposes
	rect := image.Rect(x1, y1, x2, y2)
	gocv.Rectangle(&img, rect, color.RGBA{255, 0, 0, 0}, 2)

	// Display the resulting image
	window := gocv.NewWindow("Saliency Image")
	defer window.Close()
	window.IMShow(img)
	window.WaitKey(0)

	// Print the result
	fmt.Println("Number of white pixels in ROI:", numWhitePixels)

	// Print the result
	fmt.Println("Number of white pixels in all image:", numWhitePixelsAll)

	fmt.Printf("Aoi: %.2f%%\n", (float64(numWhitePixels)/float64(numWhitePixelsAll))*100)

}
