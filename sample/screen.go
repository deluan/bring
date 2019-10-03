package main

import (
	"image"

	"github.com/faiface/pixel"
)

func (app *SampleApp) updateScreen(img image.Image) {
	winWidth := app.win.Bounds().Max.X
	winHeight := app.win.Bounds().Max.Y

	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	// Scale and center image
	scale := pixel.V(winWidth/float64(imgWidth), winHeight/float64(imgHeight))
	mat := pixel.IM
	mat = mat.ScaledXY(pixel.ZV, scale)
	mat = mat.Moved(app.win.Bounds().Center())

	// Put image in a sprite
	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())

	// Renders the sprite in the window, using scaled matrix
	sprite.Draw(app.win, mat)
}
