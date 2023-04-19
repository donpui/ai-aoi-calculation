# ai-aoi-calculation
Calculation Area of Interest in Image or Video from predictive eye-tracking saliency map.
Calculated values show how much attention is received in a certain area compared to all images.

## Info

Built with help of ChatGPT 

Requires video or images generated with Attention Insight predictive service.

## Run
Currently, all values are set in `calculate`.go`

`go run .`

## TODO

- [] Allow users to draw aoi on extracted frame
- [] Improve calculation logic to calculate averages or median in certain periods of frames