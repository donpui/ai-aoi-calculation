# ai-aoi-calculation

Calculation Area of Interest in Image or Video from predictive eye-tracking saliency map.
Calculated values show how much attention is received in a certain area compared to all images.

## Info

Built with help of ChatGPT 

Requires video or images generated with Attention Insight predictive service.

## Pre-requisite:

- Install opencv library (macOS: `brew install opencv`)

## Run
Currently, all values are set in `calculate`.go`

`go run . <heatmap-filepath> <saliency-filepath>`
or 
`./aoi-calculation <heatmap-filepath> <saliency-filepath>`

Build:
`go build .`

## TODO

- [] Improve calculation logic to calculate averages or median in certain periods of frames
- [] Add image calculation logic, extract to separate file
- [] Improve AOI numbers display (depend font size from image dimensions, display near aoi box) 
- [] Add tests
- [] Refactor code for a clear structure
- [] Add UI dialog for selecting a file. Note: https://github.com/sqweek/dialog has issues on macOS
- [] Add CI/CD for releases exec
- [] Cross-compile for AMD64/ARM64 (currently runs on ARM64)
- [] Add object detections (texts)
- [] For v2, move to cloud function or backend service with a simple web app interface
