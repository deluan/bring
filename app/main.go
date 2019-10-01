package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/deluan/bring"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/colornames"
)

const (
	mainWidth  = 1024
	mainHeight = 768
)

var stateNames = map[bring.SessionState]string{
	bring.SessionActive:    "Active",
	bring.SessionClosed:    "Closed",
	bring.SessionHandshake: "Handshake",
}

func initBring(protocol, hostname, port string) *bring.Client {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true, ForceColors: true})
	logger.SetLevel(logrus.DebugLevel)

	session, err := bring.NewSession("localhost:4822", protocol, map[string]string{
		"hostname": hostname,
		"port":     port,
		"password": "vncpassword",
		"width":    strconv.Itoa(mainWidth),
		"height":   strconv.Itoa(mainHeight),
	}, logger)
	if err != nil {
		panic(err)
	}

	client, err := bring.NewClient(session, logger)
	if err != nil {
		panic(err)
	}
	go client.Start()
	return client
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Bring it on!",
		Bounds: pixel.R(0, 0, mainWidth, mainHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	client := initBring(os.Args[1], os.Args[2], os.Args[3])

	win.Clear(colornames.Skyblue)
	win.SetCursorVisible(false)

	mat := pixel.IM
	mat = mat.Moved(win.Bounds().Center())
	//mat = mat.Scaled(win.Bounds().Center(), 900/1024)

	frames := 0
	second := time.Tick(time.Second)
	var lastRefresh int64

	var mousePos pixel.Vec
	var mouseBtns []bring.MouseButton

	for !win.Closed() {
		// Process screen updates
		img, lastUpdate := client.Canvas()
		if lastRefresh != lastUpdate {
			if img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0 {
				pic := pixel.PictureDataFromImage(img)
				sprite := pixel.NewSprite(pic, pic.Bounds())
				sprite.Draw(win, mat)
			}
			lastRefresh = lastUpdate
		}
		win.Update()

		// Handle mouse events
		newMousePos := win.MousePosition()
		newMouseBtns := mouseButtons(win)
		if mouseInWindow(win, newMousePos) &&
			(mousePos != newMousePos || !reflect.DeepEqual(mouseBtns, newMouseBtns) || changeInMouseButtons(win)) {
			y := mainHeight - mousePos.Y // OpenGL uses inverted Y
			client.MoveMouse(image.Pt(int(mousePos.X), int(y)), newMouseBtns...)
			mousePos = newMousePos
			mouseBtns = newMouseBtns
		}

		// Handle keyboard events
		pressed, released := collectKeys(win)
		for _, k := range pressed {
			client.SendKey(k, true)
		}
		for _, k := range released {
			client.SendKey(k, false)
		}

		// Measure FPS
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | %s | FPS: %d", cfg.Title, stateNames[client.State()], frames))
			frames = 0
		default:
		}
	}
}

// Rant: why pixelgl keyboard events handling is so messy?!?
func collectKeys(win *pixelgl.Window) (pressed []bring.KeyCode, released []bring.KeyCode) {
	for k, v := range keys {
		key := v
		if win.JustPressed(k) || win.Repeated(k) {
			pressed = append(pressed, key)
		}
		if win.JustReleased(k) {
			released = append(released, key)
		}
	}
	controlPressed := win.Pressed(pixelgl.KeyLeftControl) || win.Pressed(pixelgl.KeyRightControl) ||
		win.Pressed(pixelgl.KeyLeftAlt) || win.Pressed(pixelgl.KeyRightAlt)
	if controlPressed {
		shiftPressed := win.Pressed(pixelgl.KeyLeftShift) || win.Pressed(pixelgl.KeyRightShift)
		for ch := 32; ch < 127; ch++ {
			isLetter := ch >= int('A') && ch <= int('Z')
			key := ch
			if isLetter && !shiftPressed {
				key = ch + 32
			}
			if win.JustPressed(pixelgl.Button(ch)) || win.Repeated(pixelgl.Button(ch)) {
				pressed = append(pressed, bring.KeyCode{key})
			}
			if win.JustReleased(pixelgl.Button(ch)) {
				released = append(released, bring.KeyCode{key})
			}
		}
	} else {
		for _, ch := range win.Typed() {
			pressed = append(pressed, bring.KeyCode{int(ch)})
			released = append(released, bring.KeyCode{int(ch)})
		}
	}
	return
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

func mouseButtons(win *pixelgl.Window) []bring.MouseButton {
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

func mouseInWindow(win *pixelgl.Window, mousePos pixel.Vec) bool {
	return win.Bounds().Contains(mousePos)
}

func main() {
	if len(os.Args) < 4 {
		println("Usage: app <vnc|rdp> address port")
		return
	}
	initKeys()
	pixelgl.Run(run)
}
