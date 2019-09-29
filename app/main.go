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

func initBring(protocol, hostname, port string) *bring.Client {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true, ForceColors: true})
	logger.SetLevel(logrus.InfoLevel)

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

	mat := pixel.IM
	mat = mat.Moved(win.Bounds().Center())
	//mat = mat.Scaled(win.Bounds().Center(), 900/1024)

	frames := 0
	second := time.Tick(time.Second)
	var lastRefresh int64

	var mousePos pixel.Vec
	var mouseBtns []int

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
		if mouseInWindow(win) && (mousePos != newMousePos || !reflect.DeepEqual(mouseBtns, newMouseBtns) || changeInMouseButtons(win)) {
			y := mainHeight - mousePos.Y // OpenGL uses inverted Y
			client.MoveMouse(image.Pt(int(mousePos.X), int(y)), newMouseBtns...)
			mousePos = newMousePos
			mouseBtns = newMouseBtns
		}

		// Measure FPS
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
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

func mouseButtons(win *pixelgl.Window) []int {
	btnMap := map[pixelgl.Button]int{
		pixelgl.MouseButtonLeft:   bring.MouseLeft,
		pixelgl.MouseButtonRight:  bring.MouseRight,
		pixelgl.MouseButtonMiddle: bring.MouseMiddle,
	}
	var btns []int
	for p, b := range btnMap {
		if win.Pressed(p) {
			btns = append(btns, b)
		}
	}
	return btns
}

func mouseInWindow(win *pixelgl.Window) bool {
	p := win.MousePosition()
	return win.Bounds().Contains(p)
}

func main() {
	if len(os.Args) < 4 {
		println("Usage: app <vnc|rdp> address port")
		return
	}

	pixelgl.Run(run)
}
