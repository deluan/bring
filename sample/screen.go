package main

import (
	"image"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func updateScreen(win *pixelgl.Window, img image.Image) {
	winWidth := win.Bounds().Max.X
	winHeight := win.Bounds().Max.Y

	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	// Scale and center image
	scale := pixel.V(winWidth/float64(imgWidth), winHeight/float64(imgHeight))
	mat := pixel.IM
	mat = mat.ScaledXY(pixel.ZV, scale)
	mat = mat.Moved(win.Bounds().Center())

	// Put image in a sprite
	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())

	// Renders the sprite in the window, using scaled matrix
	sprite.Draw(win, mat)
}
