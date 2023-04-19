# ai-aoi-calculation
Calculation Area of Interest in Image or Video from predictive eye-tracking saliency map.
Calculated values show how much attention is received in a certain area compared to all images.

## Info

Built with help of ChatGPT 

Requires video or images generated with Attention Insight predictive service.

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
- [] Add tests
- [] Refactor code
