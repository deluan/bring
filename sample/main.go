package main

import (
	"fmt"
	"image"
	"os"
	"strconv"
	"time"

	"github.com/deluan/bring"
	"github.com/sirupsen/logrus"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowTitle   = "Bring it on!"
	defaultWidth  = 1024
	defaultHeight = 768
	guacdAddress  = "localhost:4822"
)

var stateNames = map[bring.SessionState]string{
	bring.SessionActive:    "Active",
	bring.SessionClosed:    "Closed",
	bring.SessionHandshake: "Handshake",
}

func createSDLWindow() (*sdlcanvas.Window, *canvas.Canvas) {
	win, cv, err := sdlcanvas.CreateWindow(defaultWidth, defaultHeight, windowTitle)
	if err != nil {
		panic(err)
	}
	//win.Window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	_, _ = sdl.ShowCursor(0)
	//win.Window.SetResizable(true)

	return win, cv
}

// Creates and initialize Bring's Session and Client
func createBringClient(protocol, hostname, port string) *bring.Client {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true, ForceColors: true})
	logger.SetLevel(logrus.DebugLevel)

	session, err := bring.NewSession(guacdAddress, protocol, map[string]string{
		"hostname": hostname,
		"port":     port,
		"password": "vncpassword",
		"width":    strconv.Itoa(defaultWidth),
		"height":   strconv.Itoa(defaultHeight),
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

func main() {
	win, cv := createSDLWindow()
	defer win.Destroy()
	client := createBringClient(os.Args[1], os.Args[2], os.Args[3])

	hookMouse(win, client)
	hookKeyboard(win, client)

	var lastRefresh int64
	second := time.Tick(time.Second)
	win.MainLoop(func() {
		select {
		case <-second:
			win.Window.SetTitle(fmt.Sprintf("%s | FPS: %2.0f | %s", windowTitle, win.FPS(), stateNames[client.State()]))
		default:
		}

		img, lastUpdate := client.Screen()

		imgWidth := img.Bounds().Max.X
		imgHeight := img.Bounds().Max.Y

		// If the image is empty, terminate loop
		if imgWidth == 0 || imgHeight == 0 {
			return
		}

		// If there were no changes, terminate loop
		if lastRefresh == lastUpdate {
			return
		}

		// Process screen updates
		cv.PutImageData(img.(*image.RGBA), 0, 0)
		lastRefresh = lastUpdate
	})
}
