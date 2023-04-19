package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"gocv.io/x/gocv"
)

type CloseWindowError struct{}

func (e *CloseWindowError) Error() string {
	return "closing the window"
}

type OnDrawnCallback func(d *Drawing, aoiStart, aoiEnd image.Point)

type Drawing struct {
	image    *ebiten.Image
	aoiStart image.Point
	aoiEnd   image.Point
	drawing  bool
	done     bool
	onDrawn  OnDrawnCallback
}

func (d *Drawing) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !d.drawing {
			d.drawing = true
			d.aoiStart.X, d.aoiStart.Y = ebiten.CursorPosition()
		}
		d.aoiEnd.X, d.aoiEnd.Y = ebiten.CursorPosition()
	} else if d.drawing {
		d.drawing = false
		fmt.Printf("Area of Interest: start(%d, %d) end(%d, %d)\n", d.aoiStart.X, d.aoiStart.Y, d.aoiEnd.X, d.aoiEnd.Y)
		if d.onDrawn != nil {
			d.onDrawn(d, d.aoiStart, d.aoiEnd)
		}
	}
	if d.done {
		return &CloseWindowError{}
	}
	return nil
}

// func (d *Drawing) Draw(screen *ebiten.Image) {
// 	screen.DrawImage(d.image, &ebiten.DrawImageOptions{})
// 	if d.drawing {
// 		ebitenutil.DrawRect(screen, float64(d.aoiStart.X), float64(d.aoiStart.Y), float64(d.aoiEnd.X-d.aoiStart.X), float64(d.aoiEnd.Y-d.aoiStart.Y), color.RGBA{0, 255, 0, 255})
// 	}
// }

func (d *Drawing) Draw(screen *ebiten.Image) {
	screen.DrawImage(d.image, &ebiten.DrawImageOptions{})

	if d.drawing {
		borderColor := color.RGBA{0, 255, 0, 255}
		borderWidth := 2.0 // Adjust the border width as needed

		rectWidth := float32(d.aoiEnd.X - d.aoiStart.X)
		rectHeight := float32(d.aoiEnd.Y - d.aoiStart.Y)

		if rectWidth > 0 && rectHeight > 0 {
			// Top border
			topBorder := ebiten.NewImage(int(rectWidth), int(borderWidth))
			topBorder.Fill(borderColor)
			topBorderOp := &ebiten.DrawImageOptions{}
			topBorderOp.GeoM.Translate(float64(d.aoiStart.X), float64(d.aoiStart.Y))
			screen.DrawImage(topBorder, topBorderOp)

			// Bottom border
			bottomBorder := ebiten.NewImage(int(rectWidth), int(borderWidth))
			bottomBorder.Fill(borderColor)
			bottomBorderOp := &ebiten.DrawImageOptions{}
			bottomBorderOp.GeoM.Translate(float64(d.aoiStart.X), float64(d.aoiStart.Y+int(rectHeight)-int(borderWidth)))
			screen.DrawImage(bottomBorder, bottomBorderOp)

			// Left border
			leftBorder := ebiten.NewImage(int(borderWidth), int(rectHeight))
			leftBorder.Fill(borderColor)
			leftBorderOp := &ebiten.DrawImageOptions{}
			leftBorderOp.GeoM.Translate(float64(d.aoiStart.X), float64(d.aoiStart.Y))
			screen.DrawImage(leftBorder, leftBorderOp)

			// Right border
			rightBorder := ebiten.NewImage(int(borderWidth), int(rectHeight))
			rightBorder.Fill(borderColor)
			rightBorderOp := &ebiten.DrawImageOptions{}
			rightBorderOp.GeoM.Translate(float64(d.aoiStart.X+int(rectWidth)-int(borderWidth)), float64(d.aoiStart.Y))
			screen.DrawImage(rightBorder, rightBorderOp)
		}
	}
}

func (d *Drawing) Layout(outsideWidth, outsideHeight int) (int, int) {
	return d.image.Size()
}

func extractFrame(videoFile string, framePosition float64) (image.Image, error) {
	video, err := gocv.VideoCaptureFile(videoFile)
	if err != nil {
		return nil, err
	}
	defer video.Close()

	video.Set(gocv.VideoCapturePosMsec, framePosition*1000)
	img := gocv.NewMat()
	defer img.Close()

	if ok := video.Read(&img); !ok {
		return nil, fmt.Errorf("failed to read frame from video")
	}
	if img.Empty() {
		return nil, fmt.Errorf("no frame captured")
	}

	nativeBuf, err := gocv.IMEncode(".jpg", img)
	if err != nil {
		return nil, err
	}
	buf := nativeBuf.GetBytes()
	return jpeg.Decode(bytes.NewReader(buf))

}

func draw(video string) *aoi_c {
	var wg sync.WaitGroup
	wg.Add(1)

	var aoiCoords *aoi_c
	frame, err := extractFrame(video, 1) // Extract 1 second position frame
	if err != nil {
		fmt.Printf("Error extracting frame: %v\n", err)
		return nil
	}

	ebitenImg := ebiten.NewImageFromImage(frame)

	drawing := &Drawing{
		image: ebitenImg,
		onDrawn: func(d *Drawing, aoiStart, aoiEnd image.Point) {
			fmt.Printf("Area of Interest: start(%d, %d) end(%d, %d)\n", aoiStart.X, aoiStart.Y, aoiEnd.X, aoiEnd.Y)
			aoiCoords = &aoi_c{
				x_start: aoiStart.X,
				y_start: aoiStart.Y,
				x_end:   aoiEnd.X,
				y_end:   aoiEnd.Y,
			}
			wg.Done()
			d.done = true
		},
	}

	ebiten.SetWindowTitle("Draw Area of Interest")
	ebiten.SetRunnableOnUnfocused(true)
	if err := ebiten.RunGame(drawing); err != nil && !errors.Is(err, &CloseWindowError{}) {
		fmt.Printf("Error running Ebiten game: %v\n", err)
		return nil
	}

	wg.Wait()

	return aoiCoords
}
