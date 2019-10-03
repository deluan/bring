package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"time"

	"github.com/deluan/bring"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/colornames"
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

func mainLoop(win *pixelgl.Window, client *bring.Client) {
	frames := 0
	second := time.Tick(time.Second)
	var lastRefresh int64

	for !win.Closed() {
		// Get an updated image from the Bring Client
		img, lastUpdate := client.Screen()

		imgWidth := img.Bounds().Dx()
		imgHeight := img.Bounds().Dy()

		// If the image is not empty
		if imgWidth > 0 && imgHeight > 0 {

			// Process screen updates if there were any updates in the image
			if lastRefresh != lastUpdate {
				updateScreen(win, img)
				lastRefresh = lastUpdate
			}

			// Handle mouse events
			mouseInfo := collectNewMouseInfo(win, imgWidth, imgHeight)
			if mouseInfo != nil {
				if err := client.SendMouse(mouseInfo.pos, mouseInfo.pressedButtons...); err != nil {
					fmt.Printf("Error: %s", err)
				}
			}

			// Handle keyboard events
			pressed, released := collectKeyStrokes(win)
			for _, k := range pressed {
				if err := client.SendKey(k, true); err != nil {
					fmt.Printf("Error: %s", err)
				}
			}
			for _, k := range released {
				if err := client.SendKey(k, false); err != nil {
					fmt.Printf("Error: %s", err)
				}
			}
		}

		win.Update()

		// Measure FPS and update title
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | %s | FPS: %d", windowTitle, stateNames[client.State()], frames))
			frames = 0
		default:
		}
	}
}

// Create the App's main window
func createAppWindow() *pixelgl.Window {
	cfg := pixelgl.WindowConfig{
		Title:     windowTitle,
		Bounds:    pixel.R(0, 0, defaultWidth, defaultHeight),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.Clear(colornames.Skyblue)
	win.SetCursorVisible(false)
	return win
}

// Pixel library requires the main to be run inside pixelgl.Run, to guarantee it is run in the main thread
func Main() {
	if len(os.Args) < 4 {
		println("Usage: app <vnc|rdp> address port")
		return
	}
	client := createBringClient(os.Args[1], os.Args[2], os.Args[3])
	win := createAppWindow()
	mainLoop(win, client)
}

func main() {
	pixelgl.Run(Main)
}
