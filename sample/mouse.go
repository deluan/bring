package main

import (
	"image"
	"reflect"

	"github.com/deluan/bring"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type mouseInfo struct {
	pos            image.Point
	pressedButtons []bring.MouseButton
}

func (app *SampleApp) collectNewMouseInfo(imgWidth, imgHeight int) *mouseInfo {
	newMousePos := app.win.MousePosition()
	newMouseBtns := pressedMouseButtons(app.win)

	// If mouse is inside window boundaries and anything has changed
	if app.win.MouseInsideWindow() &&
		(app.win.MousePreviousPosition() != newMousePos ||
			!reflect.DeepEqual(app.mousePreviousButtons, newMouseBtns)) {

		app.mousePreviousButtons = newMouseBtns

		winWidth := app.win.Bounds().Max.X
		winHeight := app.win.Bounds().Max.Y

		// Scale mouse position
		scale := pixel.V(float64(imgWidth)/winWidth, float64(imgHeight)/winHeight)
		newMousePos = newMousePos.ScaledXY(scale)

		// OpenGL uses inverted Y
		y := float64(imgHeight) - newMousePos.Y

		pos := image.Pt(int(newMousePos.X), int(y))
		return &mouseInfo{pos, newMouseBtns}
	}
	return nil
}

func pressedMouseButtons(win *pixelgl.Window) []bring.MouseButton {
	btnMap := map[pixelgl.Button]bring.MouseButton{
		pixelgl.MouseButtonLeft:   bring.MouseLeft,
		pixelgl.MouseButtonRight:  bring.MouseRight,
		pixelgl.MouseButtonMiddle: bring.MouseMiddle,
	}
	var btns []bring.MouseButton
	for p, b := range btnMap {
		if win.Pressed(p) {
			btns = append(btns, b)
		}
	}
	return btns
}
