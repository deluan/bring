package main

import (
	"fmt"
	"image"

	"github.com/deluan/bring"
	"github.com/tfriedel6/canvas/sdlcanvas"
)

type mouseState struct {
	client  *bring.Client
	buttons map[int]bool
	x, y    int
}

var mouseBtnMap = map[int]bring.MouseButton{
	1: bring.MouseLeft,
	2: bring.MouseMiddle,
	3: bring.MouseRight,
}

func hookMouse(win *sdlcanvas.Window, client *bring.Client) {
	ms := &mouseState{client: client, buttons: make(map[int]bool)}
	win.MouseMove = ms.mouseMove
	win.MouseDown = ms.mouseDown
	win.MouseUp = ms.mouseUp
}

func (ms *mouseState) pressedButtons() []bring.MouseButton {
	var buttons []bring.MouseButton
	for b, pressed := range ms.buttons {
		bb := mouseBtnMap[b]
		if pressed {
			buttons = append(buttons, bb)
		}
	}
	return buttons
}

func (ms *mouseState) sendMouse(x, y int) {
	ms.x, ms.y = x, y
	if err := ms.client.SendMouse(image.Pt(x, y), ms.pressedButtons()...); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func (ms *mouseState) mouseDown(button, x, y int) {
	ms.buttons[button] = true
	ms.sendMouse(x, y)
}

func (ms *mouseState) mouseUp(button, x, y int) {
	ms.buttons[button] = false
	ms.sendMouse(x, y)
}

func (ms *mouseState) mouseMove(x, y int) {
	if ms.x == x && ms.y == y {
		return
	}
	ms.sendMouse(x, y)
}
