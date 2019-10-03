package main

import (
	"image"

	"github.com/deluan/bring"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type mouseInfo struct {
	pos            image.Point
	pressedButtons []bring.MouseButton
}

func collectNewMouseInfo(win *pixelgl.Window, imgWidth, imgHeight int) *mouseInfo {
	newMousePos := win.MousePosition()
	newMouseBtns := pressedMouseButtons(win)

	// If mouse is inside window boundaries and anything has changed
	if win.MouseInsideWindow() &&
		(win.MousePreviousPosition() != newMousePos || changeInMouseButtons(win)) {

		winWidth := win.Bounds().Max.X
		winHeight := win.Bounds().Max.Y

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

func changeInMouseButtons(win *pixelgl.Window) bool {
	btns := []pixelgl.Button{
		pixelgl.MouseButtonLeft,
		pixelgl.MouseButtonRight,
		pixelgl.MouseButtonMiddle,
	}
	for _, p := range btns {
		if win.JustPressed(p) || win.JustReleased(p) {
			return true
		}
	}
	return false
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
